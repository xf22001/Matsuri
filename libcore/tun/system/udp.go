package system

import (
	"libcore/tun"
	"net"

	"github.com/v2fly/v2ray-core/v5/common/buf"
	"gvisor.dev/gvisor/pkg/tcpip"
	"gvisor.dev/gvisor/pkg/tcpip/buffer"
	"gvisor.dev/gvisor/pkg/tcpip/header"
)

func (t *SystemTun) processIPv4UDP(cache *buf.Buffer, ipHdr header.IPv4, hdr header.UDP) {
	sourceAddress := ipHdr.SourceAddress()
	destinationAddress := ipHdr.DestinationAddress()
	sourcePort := hdr.SourcePort()
	destinationPort := hdr.DestinationPort()

	source := &net.UDPAddr{
		IP:   []byte(sourceAddress),
		Port: int(sourcePort),
	}
	destination := &net.UDPAddr{
		IP:   []byte(destinationAddress),
		Port: int(destinationPort),
	}

	ipHdr.SetDestinationAddress(sourceAddress)
	hdr.SetDestinationPort(sourcePort)

	headerLength := ipHdr.HeaderLength()
	headerCache := buf.New()
	headerCache.Write(ipHdr[:headerLength+header.UDPMinimumSize])

	cache.Advance(int32(headerLength + header.UDPMinimumSize))
	t.handler.HandlePacket(&tun.UDPPacket{
		Src:       source,
		Dst:       destination,
		Data:      cache.Bytes(),
		PutHeader: headerCache.Release,
		WriteBack: func(bytes []byte, addr *net.UDPAddr) (int, error) {
			index := headerCache.Len()
			newHeader := headerCache.ExtendCopy(headerCache.Bytes())
			headerCache.Advance(index)

			defer func() {
				headerCache.Clear()
				headerCache.Resize(0, index)
			}()

			var newSourceAddress tcpip.Address
			var newSourcePort uint16

			if addr != nil {
				newSourceAddress = tcpip.Address(addr.IP)
				newSourcePort = uint16(addr.Port)
			} else {
				newSourceAddress = destinationAddress
				newSourcePort = destinationPort
			}

			newIpHdr := header.IPv4(newHeader)
			newIpHdr.SetSourceAddress(newSourceAddress)
			newIpHdr.SetTotalLength(uint16(int(headerCache.Len()) + len(bytes)))
			newIpHdr.SetChecksum(0)
			newIpHdr.SetChecksum(^newIpHdr.CalculateChecksum())

			udpHdr := header.UDP(headerCache.BytesFrom(headerCache.Len() - header.UDPMinimumSize))
			udpHdr.SetSourcePort(newSourcePort)
			udpHdr.SetLength(uint16(header.UDPMinimumSize + len(bytes)))
			udpHdr.SetChecksum(0)
			udpHdr.SetChecksum(^udpHdr.CalculateChecksum(header.Checksum(bytes, header.PseudoHeaderChecksum(header.UDPProtocolNumber, newSourceAddress, sourceAddress, uint16(header.UDPMinimumSize+len(bytes))))))

			replyVV := buffer.VectorisedView{}
			replyVV.AppendView(newHeader)
			replyVV.AppendView(bytes)

			if err := t.writeRawPacket(replyVV); err != nil {
				return 0, newError(err.String())
			}

			return len(bytes), nil
		},
	})
}

func (t *SystemTun) processIPv6UDP(cache *buf.Buffer, ipHdr header.IPv6, hdr header.UDP) {
	sourceAddress := ipHdr.SourceAddress()
	destinationAddress := ipHdr.DestinationAddress()
	sourcePort := hdr.SourcePort()
	destinationPort := hdr.DestinationPort()

	source := &net.UDPAddr{
		IP:   []byte(sourceAddress),
		Port: int(sourcePort),
	}
	destination := &net.UDPAddr{
		IP:   []byte(destinationAddress),
		Port: int(destinationPort),
	}

	ipHdr.SetDestinationAddress(sourceAddress)
	hdr.SetDestinationPort(sourcePort)

	headerLength := uint16(len(ipHdr)) - ipHdr.PayloadLength()
	headerCache := buf.New()
	headerCache.Write(ipHdr[:headerLength+header.UDPMinimumSize])

	cache.Advance(int32(headerLength + header.UDPMinimumSize))
	t.handler.HandlePacket(&tun.UDPPacket{
		Src:       source,
		Dst:       destination,
		Data:      cache.Bytes(),
		PutHeader: headerCache.Release,
		WriteBack: func(bytes []byte, addr *net.UDPAddr) (int, error) {
			index := headerCache.Len()
			newHeader := headerCache.ExtendCopy(headerCache.Bytes())
			headerCache.Advance(index)

			defer func() {
				headerCache.Clear()
				headerCache.Resize(0, index)
			}()

			var newSourceAddress tcpip.Address
			var newSourcePort uint16

			if addr != nil {
				newSourceAddress = tcpip.Address(addr.IP)
				newSourcePort = uint16(addr.Port)
			} else {
				newSourceAddress = destinationAddress
				newSourcePort = destinationPort
			}

			newIpHdr := header.IPv6(newHeader)
			newIpHdr.SetSourceAddress(newSourceAddress)
			newIpHdr.SetPayloadLength(uint16(header.UDPMinimumSize + len(bytes)))

			udpHdr := header.UDP(headerCache.BytesFrom(headerCache.Len() - header.UDPMinimumSize))
			udpHdr.SetSourcePort(newSourcePort)
			udpHdr.SetLength(uint16(header.UDPMinimumSize + len(bytes)))
			udpHdr.SetChecksum(0)
			udpHdr.SetChecksum(^udpHdr.CalculateChecksum(header.Checksum(bytes, header.PseudoHeaderChecksum(header.UDPProtocolNumber, newSourceAddress, sourceAddress, uint16(header.UDPMinimumSize+len(bytes))))))

			replyVV := buffer.VectorisedView{}
			replyVV.AppendView(headerCache.Bytes())
			replyVV.AppendView(bytes)

			if err := t.writeRawPacket(replyVV); err != nil {
				return 0, newError(err.String())
			}

			return len(bytes), nil
		},
	})
}
