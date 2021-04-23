// crawler - crawls all over the p2p network to figure out all the nodes alive
// out there.
//
// the crawler starts from a root connected node and based on the
// `local_peerlist_new` that it receives from a handshake, figures out other
// peers and from that, goes on and on to find more new peers, eventually
// getting a good grasp of all the peers in the network.
//
package crawler

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/cirocosta/go-monero/pkg/levin"
)

const (
	// CrawlerConcurrency determines how many goroutines are spawn to visit
	// peers.
	//
	CrawlerConcurrency = 250

	// CrawlerDiscoveryTimeout determines the maximum amount of time that a
	// worker trying to visit a peer should take before timing out and
	// giving up on it.
	//
	CrawlerDiscoveryTimeout = 30 * time.Second
)

type WorkerStatuses []byte

func (ws WorkerStatuses) Start(idx int) {
	ws[idx] = 1
}

func (ws WorkerStatuses) Finish(idx int) {
	ws[idx] = 0
}

func (ws WorkerStatuses) DoingWork() bool {
	for _, b := range ws {
		if b > 0 {
			return true
		}
	}

	return false
}

type VisitedPeer struct {
	Peer  *levin.Peer
	Error error
}

func (n *VisitedPeer) String() string {
	return n.Addr()
}

func (n *VisitedPeer) Addr() string {
	return n.Peer.Addr()
}

type Crawler struct {
	visited    map[string]*VisitedPeer
	notVisited map[string]*levin.Peer
	log        *log.Entry

	workerStatuses WorkerStatuses

	sync.Mutex
}

func NewCrawler() *Crawler {
	return &Crawler{
		visited:        map[string]*VisitedPeer{},
		notVisited:     map[string]*levin.Peer{},
		workerStatuses: make(WorkerStatuses, CrawlerConcurrency),

		log: log.WithFields(log.Fields{
			"component": "crawler",
		}),
	}
}

func (c *Crawler) TakeNotVisited() *levin.Peer {
	c.Lock()
	defer c.Unlock()

	for k, peer := range c.notVisited {
		delete(c.notVisited, k)
		return peer
	}

	return nil
}

func (c *Crawler) TryPutForVisit(peer *levin.Peer) {
	c.Lock()
	defer c.Unlock()

	if _, alreadyVisited := c.visited[peer.Addr()]; alreadyVisited {
		return
	}

	c.notVisited[peer.Addr()] = peer
	return
}

func (c *Crawler) MarkVisited(node *VisitedPeer) {
	c.Lock()
	defer c.Unlock()

	c.visited[node.Addr()] = node
	return
}

// Runs takes care of spinning up worker goroutines that will then do the job
// of communicating with the nodes and figuring out their peerlist.
//
func (c *Crawler) Run(ctx context.Context) (map[string]*VisitedPeer, error) {
	var (
		peersToVisitC  = make(chan *levin.Peer, 0)
		newPeersFoundC = make(chan *levin.Peer, 0)
		peersVisitedC  = make(chan *VisitedPeer, 0)
	)

	for workerIndex := 0; workerIndex < CrawlerConcurrency; workerIndex++ {
		go c.Worker(workerIndex, peersToVisitC, newPeersFoundC, peersVisitedC)
	}

	go func() {
		for peerFound := range newPeersFoundC {
			c.TryPutForVisit(peerFound)
			c.log.WithField("peerfound", peerFound).Debug("peer found, putting for visit")
		}
	}()

	go func() {
		for peerVisited := range peersVisitedC {
			c.MarkVisited(peerVisited)
			c.log.WithField("peervisited", peerVisited).Info("peer visited, marking visited")
		}
	}()

	for {
		peer := c.TakeNotVisited()
		if peer == nil {
			c.log.Info("no one to visit, sleeping a bit")
			time.Sleep(500 * time.Millisecond)
			continue
		}

		peersToVisitC <- peer
	}

	close(peersToVisitC)
	close(newPeersFoundC)
	close(peersVisitedC)

	return nil, nil
}

// Worker takes care of consuming a peer at a time and discovering their peers.
//
// Essentially: peerC -> Worker(peer) -> peer-peers...
//
func (c *Crawler) Worker(
	workerIndex int,
	peersToVisitC chan *levin.Peer,
	newPeersFoundC chan *levin.Peer,
	peersVisitedC chan *VisitedPeer,
) {
	logger := c.log.WithField("worker-idx", workerIndex)

	for peerToVisit := range peersToVisitC {
		l := logger.WithField("peer", peerToVisit)

		func() {
			l.Debug("visiting")
			defer l.Debug("visited")

			c.workerStatuses.Start(workerIndex)
			defer c.workerStatuses.Finish(workerIndex)

			visitedPeer := &VisitedPeer{
				Peer:  peerToVisit,
				Error: nil,
			}

			peersFound, err := c.Discover(context.Background(), peerToVisit)
			if err != nil {
				l.Error("errored discovering peers", err)
				visitedPeer.Error = err
				peersVisitedC <- visitedPeer
				return
			}

			l.WithField("found", len(peersFound)).Info("peers found")

			for _, peerFound := range peersFound {
				newPeersFoundC <- peerFound

			}

			peersVisitedC <- visitedPeer
		}()

	}

	logger.Info("peers to visit closed - bailing out")
}

func (c *Crawler) Discover(ctx context.Context, peer *levin.Peer) (map[string]*levin.Peer, error) {
	client, err := levin.NewClient(peer.Addr())
	if err != nil {
		return nil, fmt.Errorf("new client: %w", err)
	}

	defer client.Close()

	pl, err := client.Handshake(ctx)
	if err != nil {
		return nil, fmt.Errorf("handshake: %w", err)
	}

	return pl.Peers, nil
}

func IsOkError(err error) bool {
	if err, ok := err.(net.Error); ok && err.Timeout() {
		return true
	}

	return false
}
