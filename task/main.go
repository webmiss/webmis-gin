package task

type Main struct{ Base }

/* 首页 */
func (r Main) New() {
	r.Print("Cli")
}
