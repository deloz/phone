package phone

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
)

const (
	PHONE_INDEX_LENGTH = 9
)

var (
	notFoundError           = errors.New("没有找到")
	invalidPhoneNumberError = errors.New("需要手机号前七位")
	PhoneTypes              = map[int]string{
		1: "移动",
		2: "联通",
		3: "电信",
		4: "电信虚拟运营商",
		5: "联通虚拟运营商",
		6: "移动虚拟运营商",
	}
)

type PhoneInfo struct {
	Phone            string
	Province         string
	City             string
	ZipCode          int
	AreaCode         int
	PhoneType        string
	PhoneRecordCount int64
}

type Phone struct {
	Version          string
	Phone            string
	PhoneType        int
	PhoneNum         int64
	firstIndexOffset int64
	phoneRecordCount int64
	fileSize         int64
	file             *os.File
	recordContent    []byte
}

func (p *Phone) search() (info *PhoneInfo, err error) {
	_, filename, _, _ := runtime.Caller(1)
	currentDir := path.Dir(filename)

	filePath := path.Join(currentDir, "phone.dat")
	f, err := os.Open(filePath)
	if err != nil {
		return info, err
	}
	p.file = f
	defer f.Close()

	fi, err := os.Stat(filePath)
	if err != nil {
		return info, err
	}
	p.fileSize = fi.Size()

	buf, err := fileReadOffset(p.file, 0, 4)
	if err != nil {
		return info, err
	}
	p.Version = string(buf)

	buf, err = fileReadOffset(p.file, 4, 4)
	if err != nil {
		return info, err
	}
	firstIndexOffset, err := bytesToInt64(buf)
	if err != nil {
		return info, err
	}
	p.firstIndexOffset = firstIndexOffset

	p.phoneRecordCount = (p.fileSize - p.firstIndexOffset) / PHONE_INDEX_LENGTH

	return p.binary_search()
}

//二分查找
func (p *Phone) binary_search() (info *PhoneInfo, err error) {
	left := int64(0)
	right := p.phoneRecordCount - 1
	var middle, currentOffset int64

	for left <= right {
		middle = left + ((right - left) >> 1)
		currentOffset = p.firstIndexOffset + middle*PHONE_INDEX_LENGTH
		if currentOffset >= p.fileSize {
			return info, notFoundError
		}

		//索引区的每条记录的总长度为9字节, 每条记录的格式为: <手机号前七位(4字节)><记录区的偏移(4字节)><卡类型(长1字节)>
		indexRecord := make([]byte, PHONE_INDEX_LENGTH)
		_, err := p.file.ReadAt(indexRecord, currentOffset)
		if err == io.EOF || err != nil {
			return info, err
		}

		currentPhone, err := bytesToInt64(indexRecord[:4])
		if err != nil {
			return info, err
		}

		if currentPhone > p.PhoneNum {
			right = middle - 1
		} else if currentPhone < p.PhoneNum {
			left = middle + 1
		} else {
			//最后一位是手机的类型
			buf := bytes.NewReader(indexRecord[8:])
			var pt byte
			err := binary.Read(buf, binary.LittleEndian, &pt)
			if err != nil {
				return info, err
			}
			p.PhoneType = int(pt)

			//记录区的位移
			recordOffset, err := bytesToInt64(indexRecord[4:8])
			if err != nil {
				return info, err
			}

			for {
				buf, err := fileReadOffset(p.file, recordOffset, 1)
				if err != nil {
					return info, err
				}
				p.recordContent = append(p.recordContent, buf...)
				if string(buf[0:1]) == "\x00" {
					break
				}
				recordOffset++
			}

			return p.format_phone_info(), nil
		}

	}

	return info, notFoundError
}

func (p *Phone) format_phone_info() *PhoneInfo {
	info := strings.Split(string(p.recordContent), "|")
	zipCode, _ := strconv.Atoi(info[2])
	areaCode, _ := strconv.Atoi(info[3])

	return &PhoneInfo{
		Phone:            p.Phone,
		Province:         info[0],
		City:             info[1],
		ZipCode:          zipCode,
		AreaCode:         areaCode,
		PhoneType:        PhoneTypes[p.PhoneType],
		PhoneRecordCount: p.phoneRecordCount,
	}
}

func fileReadOffset(f *os.File, offset int64, length int) (b []byte, err error) {
	bs := make([]byte, length)
	_, err = f.ReadAt(bs, offset)
	if err != nil {
		return b, err
	}
	if err == io.EOF {
		return b, err
	}

	return bs, nil
}

func bytesToInt64(bs []byte) (int64, error) {
	buf := bytes.NewReader(bs)
	var num uint32
	err := binary.Read(buf, binary.LittleEndian, &num)
	if err != nil {
		return 0, err
	}

	return int64(num), nil
}

func Find(fullPhone string) (info *PhoneInfo, err error) {
	phoneNum, err := validatePhone(fullPhone)
	if err != nil {
		return info, err
	}

	p := Phone{
		Phone:    fullPhone,
		PhoneNum: int64(phoneNum),
	}

	return p.search()
}

func validatePhone(phone string) (int, error) {
	phone = strings.Trim(phone[0:7], " ")
	phoneNum, err := strconv.Atoi(phone)
	if err != nil {
		return 0, err
	}
	length := len(phone)
	if length <= 11 && length >= 7 {
		return phoneNum, nil
	}

	return 0, errors.New("invalid phone number")
}
