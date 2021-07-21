省略 if err != nil 的编写

并自动带上错误位置和函数调用信息

    package main

    func fn1() (i int, err autoerr) {
        _, err := strconv.Atoi("abc")
        return
    }

=>

    package main

    func fn1() (i int, err autoerr) {
        _, err := strconv.Atoi("abc")
        if err != nil {
            err = errors.Wrap(err, "fn `fn1` failed at apple.go:4 when invoking `strconv.Atoi`"+fmt.Sprintf(" with arg0 = %v; ", "abc"))
            return
        }
        return
    }