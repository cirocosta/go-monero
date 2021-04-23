package levin

import (
	"fmt"
	"net"
)

type LocalPeerList struct {
	Peers map[string]*Peer
}

func (l *LocalPeerList) GetPeers() map[string]*Peer {
	return l.Peers
}

type Peer struct {
	Ip   string
	Port uint16

	Id      string
	RPCPort uint16

	CurrentHeight uint64
	TopVersion    uint8
}

func (p Peer) Addr() string {
	return fmt.Sprintf("%s:%d", p.Ip, p.Port)
}

func (p Peer) String() string {
	return p.Addr()
}

// TODO less panic'ing
func NewLocalPeerListFromEntries(entries Entries) LocalPeerList {
	peers := map[string]*Peer{}

	for _, entry := range entries {
		if entry.Name != "local_peerlist_new" {
			continue
		}

		peerList := entry.Entries()

		for _, peer := range peerList {
			peerListAdr := peer.Entries()

			for _, adr := range peerListAdr {
				if adr.Name != "adr" {
					continue
				}

				addr := adr.Entries()

				for _, addrField := range addr {

					if addrField.Name != "addr" {
						continue
					}

					fields := addrField.Entries()

					var ip string
					var port uint16

					for _, field := range fields {
						if field.Name == "m_ip" {
							ip = ipzify(field.Uint32())
						}

						if field.Name == "m_port" {
							port = field.Uint16()
						}

						if field.Name == "addr" {
							ip = net.IP([]byte(field.String())).String()
						}
					}

					if ip != "" && port != 0 {

						peer := &Peer{
							Ip:   ip,
							Port: port,
						}

						peers[peer.Addr()] = peer
					}
				}
			}
		}
	}

	return LocalPeerList{
		Peers: peers,
	}

}

func ipzify(ip uint32) string {
	result := make(net.IP, 4)

	result[0] = byte(ip)
	result[1] = byte(ip >> 8)
	result[2] = byte(ip >> 16)
	result[3] = byte(ip >> 24)

	return result.String()
}
