# 代码结构

errgen 、recovergen 这些组件主要由两部分组成

visitor 和 template

visitor 负责扫描 AST 发现要生成代码的地方 / 代码插入点

template 负责生成要插入的代码

# 组件

errgen 、 recovergen 、 wiregen 等是 g 的"组件"， 它们都有一个函数签名为
    func(files []*dst.File, filenames []string, pkg *decorator.Package) (err error)
的函数作为组件的执行入口