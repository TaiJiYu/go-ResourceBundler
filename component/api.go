package component

type UnpackOption struct {
	FqmFilePath string // 需要读取的fqm文件路径
	UseCache    bool   // 是否使用缓存
	SecretKey   []byte // 密码
}

type PackOption struct {
	SecretKey   []byte   // 密码
	Name        string   // 输出文件名
	ResourceDir string   // 打包资源所在的绝对路径
	OutcomeDir  string   // 输出文件的绝对路径
	IgnoreDir   []string // 忽略的文件夹
}

type packer = *packageByte

// pack
func NewPacker(options ...PackOption) packer     { return newPackage(options...) }
func (p packer) Add(key string, data []byte)     { p.add(key, data) }
func (p packer) AddResourceDir(dir string)       { p.addResourceDir(dir) }
func (p packer) AddFile(key string, path string) { p.addFile(key, path) }
func (p packer) PackResource()                   { p.packResource() }
func (p packer) Save() error                     { return p.save() }

// unpack
func UnpackerInit(o ...UnpackOption) error { return unpackInit(o...) }
func Key(k string) []byte                  { return key(k) }
func Close()                               { close() }
func Show()                                { show() }
