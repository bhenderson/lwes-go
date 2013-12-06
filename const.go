package lwes

const (
	/*
	   from lwes c library:
	   maximum datagram size for UDP is 64K minus IP layer overhead which is
	   20 bytes for IP header, and 8 bytes for UDP header, so this value
	   should be

	   65535 - 28 = 65507
	*/
	MAX_MSG_SIZE          = 65507
	MAX_SHORT_STRING_SIZE = 255
)

const (
	// type map
	_ byte = iota
	LWES_U_INT_16_TOKEN
	LWES_INT_16_TOKEN
	LWES_U_INT_32_TOKEN
	LWES_INT_32_TOKEN
	LWES_STRING_TOKEN
	LWES_IP_ADDR_TOKEN
	LWES_INT_64_TOKEN
	LWES_U_INT_64_TOKEN
	LWES_BOOLEAN_TOKEN
)
