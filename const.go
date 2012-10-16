package lwes

const (
    /*
        from lwes c library:
        maximum datagram size for UDP is 64K minus IP layer overhead which is
        20 bytes for IP header, and 8 bytes for UDP header, so this value
        should be

        65535 - 28 = 65507
     */
    MAX_MSG_SIZE = 65507
    MAX_SHORT_STRING_SIZE = 255
)

