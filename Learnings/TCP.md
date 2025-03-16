2 computers talking to each other should know some required data:
1. message length
2. Actual data
3. Sender address
4. Receiver Address

Ground rules to communicate successfully -
1. Network Protocols

TCP/IP is one such Network Protocols.

TCP breaks message into smaller parts calles segments, and then give them to IP(routing), IP routes it to the destinations.

Transmission Control Protocol - highly reliable - slow
User Datagram Protocol - less reliable - faster

Internet Protocol (IP)

TCP/IP protocol - because tcp was famoous - TCP and IP are two differenrt layers in the 7 OSI layers, which one I can check later.


OSI - Open system Interconnection


okay no, TCP/IP protocol is for the whole communication layers but we just call it tcp ip.


application layer generates the message
each layer encapsulates the message with necessary addons for communication

Application Layer -> Application message
Transport -> TCP Segment, UDP Datagram
Network -> IP Packet
Data Link -> Ethernet frame
Physical Layer  -> Ethernet Frame


*Physical Layer* - actual communication takes place - 0's 1's
converts the binary sequence into signals and transmits them via a medium
Ethernet Protocol is widely used at this layer


*Data Link Layer* - Ethernet Frame 
The Data link layer is divided into 2 sublayers 
1. Medium Access Control - MAC Sub layer
2. Logical Link Control - LLC Sub Layer

MAC is responsible for data encapsulation and accessing the media

### Data Encapsulation
- Adds header and trailer to the Packet 
- header contains MAC sender and receiver's address - MAC address is a unique address embedded in the device by the device manufacturer 
- trailer contains 4 bytes of code for error detection

### Accessing the medio
- Access method is used called CSMA/CD (Carrier Sense Multiple Access/Collision Detection)
- Each computer listens to the cable before sending data to the network, if the network is clear then only as scomputer tries to transmit
- If two computer tries to send data at the same time, there will be a collision and the transmission stops and both the computers try after sometime
- Delay caused by collision and retransmission is ver small and does not affect the speed of transmission of data flow

LLC is responsible for flow control and error control
- it restricts the amount of data that sender can send without overwhelming the receiver
- the receiver might have limited processing speed, ifg these limits are exceeded then the incoming data might be lost.
- receiver should inform the sendeer to slow down the transmission rate 
- flow controls restricts the number of frames that sender can send withjout overwhelming the receiver
- Error control refers to error detection and retransmission. this is done usinfg the error bytes in the trailer.
- the retransmission is done using ARQ - automatic retransmission request
- receiver send an ack if the frame is received. so if an ACK is not received by the sender, the sender sends the packet again. This is called ARQ.
- LLC can also resize the IP packets received from the network layer to fit them into the data link layer.

*Transport Layer* - This usually proviodes flow control, error control and sizing of packets so LLC layer is seldom used.


Application |  
Transport   |----> these layers are implemented in software programs in the computer OS 
Network     |

Transport Layer passes the TCP segment/UDP datagram to the network layer.

*Network layer* adds logica/IP address to the packet in the header and makes it IP packet and later routes them to the other networks.
The network Layer also determines the best path for the data delivery.
IP is the single protocol for this layer. Every computer has a unique IP address. 
IP protocol is unreliable and does not guarantee delivery not it checks for the errors. So this needs to be done in the Transport Layer.

*Transport Layer* - receives message from the aaplication layer - TCP or UDp protocol is selected
- TCP supports segmentation - 3 phases - connection establishment, data transfer, connection termination
- in the connection establishment phase, the sender sends a packet to receiver requesting an ACK. if S receives an ACK, it'll again send an ACK to R completion the phase sucessfully. 3 way tcp connection handshake
- in the data transfer phase, error free data transfer, ordered(sequence numbers are sent), retransmission of lost segments(ACK is sent for each TCP segment otherwise S sends the segment again, the duplication is solved via sequencing), congestion throttling(times is manipulated accordingly)
- in connection termination phase, the sender sends a finished message 4 way handshake 
- UDP does not support segmentation - so the sender should send data in smaller size - unreliable - lacks error checking 

*Application Layer* - user application, HTTP/HTTPS protocol, SMTP for email




### Telnet - Teletype network protocol
A network protocol that lets you communicate with the remote computers. telnet
