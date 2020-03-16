package x

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func main() {

	conn, err := net.DialTimeout("tcp", "sqweeb.net:25565", time.Duration(3)*time.Second)
	if err != nil {
		return
	}
	defer conn.Close()

	conn.SetWriteDeadline(time.Now().Add(time.Second * 5))
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))

	buf := bytes.NewBuffer([]byte{})

	// All data sent over the network (except for VarInt and VarLong) is big-endian
	err = binary.Write(buf, binary.LittleEndian, int16(-1)) // protocol version
	if err != nil {
		fmt.Println(err)
		return
	}
	// TODO: Toggle this based on the whether server has plugins or not - Spigot doesn't seem to want it.
	// err = binary.Write(buf, binary.BigEndian, []byte("sqweeb.net:25565")) // server address
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	err = binary.Write(buf, binary.BigEndian, int16(25565)) // server port
	if err != nil {
		fmt.Println(err)
		return
	}
	err = binary.Write(buf, binary.LittleEndian, int8(1)) // next state (1=status,2=login)
	if err != nil {
		fmt.Println(err)
		return
	}

	handshake := buf.Bytes()
	fmt.Println(len(handshake))
	buf.Reset()

	// handshake
	if err := writePacket(conn, 0, handshake); err != nil {
		fmt.Println(err)
		return
	}

	err = binary.Write(buf, binary.LittleEndian, int8(1)) // protocol version
	if err != nil {
		fmt.Println(err)
		return
	}
	err = binary.Write(buf, binary.BigEndian, int16(25565)) // server port
	if err != nil {
		fmt.Println(err)
		return
	}

	status := buf.Bytes()
	buf.Reset()

	// status
	if err := writePacket(conn, 0, status); err != nil {
		fmt.Println(err)
		return
	}

	// ================
	// Reading Response

	readBuf := make([]byte, 4096)
	r, err := conn.Read(readBuf)
	if err != nil {
		return
	}

	fmt.Println("Response length: ", r)
	fmt.Println(string(readBuf))

	return

}

func writePacket(c net.Conn, packetType int8, packetData []byte) error {
	var err error
	pt := make([]byte, 1)
	pt[0] = byte(packetType)
	pl := len(pt) + len(packetData)

	fmt.Printf("Outgoing packet size: %v byte(s)\n", pl)
	err = binary.Write(c, binary.LittleEndian, int64(pl))
	if err != nil {
		return err
	}

	fmt.Printf("Outgoing packet type: %v\n", pt)
	err = binary.Write(c, binary.LittleEndian, pt)
	if err != nil {
		return err
	}

	fmt.Printf("Outgoing packet data: %v\n", packetData)
	return binary.Write(c, binary.BigEndian, packetData)
}
