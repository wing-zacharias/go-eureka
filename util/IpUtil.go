package util

/**
  @author: wing
  @date: 2020/9/4
  @comment:
**/
import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

/**
* @author: wing
* @time: 2020/9/4 9:52
* @param:
* @return:
* @comment: ipaddress detail entity
**/
type IpAddress struct {
	Name string
	IP   string
	MASK int
	MAC  string
}

/**
* @author: wing
* @time: 2020/9/4 9:52
* @param:
* @return:
* @comment: get local hostname
**/
func GetHostname() string {
	hostName, _ := os.Hostname()
	return hostName
}

/**
* @author: wing
* @time: 2020/9/4 12:29
* @param:
* @return:
* @comment: ipv4 check
**/
func Ipv4Check(ipv4Str string) bool {
	ip := net.ParseIP(ipv4Str)
	if ip != nil {
		return true
	}
	return false
}

/**
* @author: wing
* @time: 2020/9/4 9:52
* @param:
* @return:
* @comment: get local ip mask.may be not accurate if net interface name not 'ens,eth,以太网'
**/
func GetLocalIp() (string, int) {
	var ip string
	var mask int
	for _, ipAddr := range GetLocalIpList() {
		if strings.Contains(ipAddr.Name, "以太网") || strings.Contains(ipAddr.Name, "ens") || strings.Contains(ipAddr.Name, "eth") {
			ip = ipAddr.IP
			mask = ipAddr.MASK
			break
		}
	}
	return ip, mask
}

/**
* @author: wing
* @time: 2020/9/4 12:31
* @param:
* @return: ip,length of mask
* @comment: cidr address disassemble
**/
func CidrDisassemble(cidrAddress string) (string, int, error) {
	parts := strings.Split(cidrAddress, "/")
	if len(parts) == 2 {
		ip := net.ParseIP(parts[0])
		maskLen, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", 0, err
		}
		if ip != nil && maskLen > 0 && maskLen < 32 {
			return parts[0], maskLen, nil
		}
	}
	return "", 0, errors.New("Input error! ")
}

/**
* @author: wing
* @time: 2020/9/4 9:55
* @param:
* @return: IpAddress
* @comment: get local ipaddress list
**/
func GetLocalIpList() []IpAddress {
	var addresses []IpAddress
	iFaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	for _, iFace := range iFaces {
		var address IpAddress
		if (iFace.Flags&net.FlagUp == 0) || (iFace.Flags&net.FlagLoopback != 0) {
			continue
		}
		addrs, err := iFace.Addrs()
		if err != nil {
			continue
		}
		address.Name = iFace.Name
		address.MAC = iFace.HardwareAddr.String()
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				{
					if v.IP == nil || v.IP.IsLoopback() {
						continue
					}
					if v.IP.To4() != nil {
						address.IP = v.IP.String()
						maskHex, _ := strconv.ParseUint(v.Mask.String(), 16, 32)
						maskBin := strconv.FormatInt(int64(maskHex), 2)
						maskLen, _ := MaskBinToMaskLen(maskBin)
						address.MASK = maskLen
					}
				}
			}
		}
		addresses = append(addresses, address)
	}
	return addresses
}

/**
* @author: wing
* @time: 2020/9/4 9:56
* @param:
* @return:
* @comment: get network flag
**/
func GetNetFlag(localIp string, maskLen int) string {
	binIp, _ := IPv4ToBin(localIp)
	pb, _ := strconv.ParseUint(binIp, 2, 32)
	maskBin, _ := MaskLenToMaskBin(maskLen)
	mb, _ := strconv.ParseUint(maskBin, 2, 32)
	r := pb & mb
	rs := strconv.FormatInt(int64(r), 16)
	res, _ := HexIToPv4(rs)
	return res
}

/**
* @author: wing
* @time: 2020/9/4 9:56
* @param:
* @return:
* @comment: get all ips in special net
**/
func GetNetIpList(localIp string, maskLen int) []string {
	var ips []string
	binIp, _ := IPv4ToBin(localIp)
	pb, _ := strconv.ParseUint(binIp, 2, 32)
	maskBin, _ := MaskLenToMaskBin(maskLen)
	mb, _ := strconv.ParseUint(maskBin, 2, 32)
	netFlag := pb & mb
	count := int(^uint32(mb))
	for i := 1; i < count; i++ {
		rs := strconv.FormatInt(int64(uint32(netFlag)+uint32(i)), 16)
		apd := 8 - len(rs)
		for i := 0; i < apd; i++ {
			rs = "0" + rs
		}
		ip, _ := HexIToPv4(rs)
		ips = append(ips, ip)
	}
	return ips
}

/**
* @author: wing
* @time: 2020/9/4 9:57
* @param:
* @return:
* @comment: mask cidr code to bin string
**/
func MaskLenToMaskBin(maskLen int) (string, error) {
	if maskLen < 0 || maskLen > 32 {
		return "", errors.New("Input error! ")
	}
	bMask := ^uint32(0) << uint(32-maskLen)
	maskBin := strconv.FormatInt(int64(bMask), 2)
	return maskBin, nil
}

/**
* @author: wing
* @time: 2020/9/4 9:58
* @param:
* @return:
* @comment: mask bin string to cidr code
**/
func MaskBinToMaskLen(maskBin string) (int, error) {
	if len(maskBin) != 32 {
		return 0, errors.New("Input error! ")
	}
	return len(strings.Replace(maskBin, "0", "", -1)), nil
}

/**
* @author: wing
* @time: 2020/9/4 9:59
* @param:
* @return:
* @comment: msak ipv4 to cidr code
**/
func MaskToMaskLen(mask string) int {
	maskBin, _ := IPv4ToBin(mask)
	maskLen, _ := MaskBinToMaskLen(maskBin)
	return maskLen
}

/**
* @author: wing
* @time: 2020/9/4 9:59
* @param:
* @return:
* @comment: mask cidr code to ipv4
**/
func MaskLenToMask(maskLen int) string {
	maskBin, _ := MaskLenToMaskBin(maskLen)
	mask, _ := BinToIPv4(maskBin)
	return mask
}

/**
* @author: wing
* @time: 2020/9/4 10:00
* @param:
* @return:
* @comment: mask ipv4 to hex
**/
func IPv4ToHex(ipv4 string) (string, error) {
	var res string
	if ip := net.ParseIP(ipv4); ip == nil {
		return "", errors.New("Input error! ")
	}
	ps := strings.Split(ipv4, ".")
	for _, p := range ps {
		pi, _ := strconv.Atoi(p)
		ph := strconv.FormatInt(int64(pi), 16)
		apd := 2 - len([]rune(ph))
		for i := 0; i < apd; i++ {
			ph = "0" + ph
		}
		res += ph
	}
	return res, nil
}

/**
* @author: wing
* @time: 2020/9/4 10:00
* @param:
* @return:
* @comment: mask hex to ipv4
**/
func HexIToPv4(hexIp string) (string, error) {
	//var res string
	if len([]rune(hexIp)) != 8 {
		return "", errors.New("Input error! ")
	}
	//ph, _ := hex.DecodeString(hexIp)
	//for idx, pp := range ph {
	//	if idx != 0 {
	//		res += "."
	//	}
	//	res += fmt.Sprintf("%d", pp)
	//}
	//return res, nil
	bIp, _ := strconv.ParseUint(hexIp, 16, 32)
	partS1 := uint8(uint32(bIp) >> 24)
	partS2 := uint8(uint32(bIp) >> 16)
	partS3 := uint8(uint32(bIp) >> 8)
	partS4 := uint8(uint32(bIp) & uint32(255))
	ip := fmt.Sprint(partS1) + "." + fmt.Sprint(partS2) + "." + fmt.Sprint(partS3) + "." + fmt.Sprint(partS4)
	return ip, nil
}

/**
* @author: wing
* @time: 2020/9/4 10:00
* @param:
* @return:
* @comment: ipv4 to bin string
**/
func IPv4ToBin(ipv4 string) (string, error) {
	var res string
	if ip := net.ParseIP(ipv4); ip == nil {
		return "", errors.New("Input error! ")
	}
	ps := strings.Split(ipv4, ".")
	for _, p := range ps {
		pi, _ := strconv.Atoi(p)
		ph := strconv.FormatInt(int64(pi), 2)
		apd := 8 - len([]rune(ph))
		for i := 0; i < apd; i++ {
			ph = "0" + ph
		}
		res += ph
	}
	return res, nil
}

/**
* @author: wing
* @time: 2020/9/4 10:00
* @param:
* @return:
* @comment: bin string to ipv4
**/
func BinToIPv4(binIp string) (string, error) {
	if len([]rune(binIp)) != 32 {
		return "", errors.New("Input error! ")
	}
	pb, _ := strconv.ParseUint(binIp, 2, 32)
	ph := strconv.FormatInt(int64(pb), 16)
	res, _ := HexIToPv4(ph)
	return res, nil
}
