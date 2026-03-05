



package unix

import (
	"errors"
	"fmt"
	"strconv"
)








func Pledge(promises, execpromises string) error {
	if err := pledgeAvailable(); err != nil {
		return err
	}

	pptr, err := BytePtrFromString(promises)
	if err != nil {
		return err
	}

	exptr, err := BytePtrFromString(execpromises)
	if err != nil {
		return err
	}

	return pledge(pptr, exptr)
}






func PledgePromises(promises string) error {
	if err := pledgeAvailable(); err != nil {
		return err
	}

	pptr, err := BytePtrFromString(promises)
	if err != nil {
		return err
	}

	return pledge(pptr, nil)
}






func PledgeExecpromises(execpromises string) error {
	if err := pledgeAvailable(); err != nil {
		return err
	}

	exptr, err := BytePtrFromString(execpromises)
	if err != nil {
		return err
	}

	return pledge(nil, exptr)
}


func majmin() (major int, minor int, err error) {
	var v Utsname
	err = Uname(&v)
	if err != nil {
		return
	}

	major, err = strconv.Atoi(string(v.Release[0]))
	if err != nil {
		err = errors.New("cannot parse major version number returned by uname")
		return
	}

	minor, err = strconv.Atoi(string(v.Release[2]))
	if err != nil {
		err = errors.New("cannot parse minor version number returned by uname")
		return
	}

	return
}



func pledgeAvailable() error {
	maj, min, err := majmin()
	if err != nil {
		return err
	}

	
	if maj < 6 || (maj == 6 && min <= 3) {
		return fmt.Errorf("cannot call Pledge on OpenBSD %d.%d", maj, min)
	}

	return nil
}
