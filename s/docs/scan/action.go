// author: asydevc <asydev@163.com>
// date: 2021-12-26

package scan

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type (
	// Action
	// API接口.
	Action interface {
		GetMethod() string

		// GetRoute
		// 读取路由.
		//
		// return "/"
		// return "/ping"
		GetRoute() string

		// GetRouteLink
		// 读取路由链接.
		//
		// return "/get.md"
		// return "/ping/get.md"
		GetRouteLink() string

		// GetTitle
		// 读取API标题.
		GetTitle() string

		Ignore() bool

		// Markdown
		// 导出Markdown文档.
		Markdown() error

		// Postman
		// 导出POSTMAN接口.
		Postman() interface{}

		// SetCommentBlock
		// 设置注释实例.
		SetCommentBlock(cb CommentBlock) Action
		SetSource(path string, line int) Action
	}

	// API结构体.
	action struct {
		// 隶属控制器.
		controller Controller

		// 接口参数.
		//
		// name=PostList
		// method=POST
		// uri=/list
		name, method, uri string

		// 行号.
		// line=102
		sourcePath string
		sourceLine int

		description, title, version string
		ignore                      bool
		request, response, sdk      string
	}
)

// NewAction
// 构造API接口实例.
func NewAction(controller Controller, name string) Action {
	o := &action{controller: controller, name: name}
	o.ignore = false
	o.title = name
	o.version = "0.0"
	o.init()
	return o
}

func (o *action) GetMethod() string { return strings.ToUpper(o.method) }

func (o *action) GetRoute() string {
	s := strings.TrimSuffix(fmt.Sprintf("%s%s", o.controller.GetRoutePrefix(), o.uri),
		"/",
	)
	if s == "" {
		s = "/"
	}
	return s
}

func (o *action) GetRouteLink() string {
	s := o.GetRoute()
	if !strings.HasSuffix(s, "/") {
		s += "/"
	}
	s = fmt.Sprintf("%s%s.md", s, strings.ToLower(o.method))
	return s
}

func (o *action) GetTitle() string { return o.title }

func (o *action) Ignore() bool { return o.ignore }

func (o *action) Markdown() (err error) {
	// 1. 忽略.
	if o.ignore {
		return
	}

	// 2. 内容.
	var text = ""
	if text, err = o.render(); err != nil {
		return
	}

	// 3. 本地.
	if o.controller.GetScanner().IsSaveEnable() {
		if err = o.controller.GetScanner().Save(
			fmt.Sprintf(
				"%s%s",
				o.controller.GetScanner().GetDocsPath(),
				o.GetRouteLink(),
			), text,
		); err != nil {
			return err
		}
	}

	// 4. 上传.
	if o.controller.GetScanner().IsUploadEnable() {
		if err = o.controller.GetScanner().Upload(o.GetRouteLink(), text); err != nil {
			return err
		}
	}

	return
}

func (o *action) Postman() interface{} {
	var (
		request = map[string]interface{}{
			"method": o.GetMethod(),
			"url": map[string]interface{}{
				"raw":  fmt.Sprintf(`{{protocol}}://%s.{{domain}}%s`, o.controller.GetScanner().GetDomainPrefix(), o.GetRoute()),
				"host": []string{o.controller.GetScanner().GetDomainPrefix(), "{{domain}}"},
				"path": strings.Split(strings.TrimPrefix(o.GetRoute(), "/"), "/"),
			},
		}
		response = make([]interface{}, 0)
	)

	// 入参.
	if o.request != "" {
		if x := o.controller.GetScanner().GetPayload(o.request); x != nil {
			if s := x.Postman(); s != "" {
				request["body"] = map[string]interface{}{
					"mode": "raw", "raw": s, "options": map[string]interface{}{
						"raw": map[string]interface{}{
							"language": "json",
						},
					},
				}
			}
		}
	}

	// 出参.
	if o.request != "" {
		if x := o.controller.GetScanner().GetPayload(o.response); x != nil {
			if s := x.Postman(); s != "" {
				response = append(response, map[string]interface{}{
					"name":   "Result",
					"body":   s,
					"status": 200,
				})
			}
		}
	}

	// n. 结果.
	return map[string]interface{}{
		"name":        o.GetTitle(),
		"description": o.description,
		"request":     request,
		"response":    response,
	}
}

func (o *action) Postman2() interface{} {
	// 入参.
	request := make(map[string]interface{})
	if o.request != "" {
		if x := o.controller.GetScanner().GetPayload(o.request); x != nil {
			if s := x.Postman(); s != "" {
				request = map[string]interface{}{
					"mode": "raw", "raw": s, "options": map[string]interface{}{
						"raw": map[string]interface{}{
							"language": "json",
						},
					},
				}
			}
		}
	}

	// 出参.
	response := make(map[string]interface{})
	if o.request != "" {
		if x := o.controller.GetScanner().GetPayload(o.response); x != nil {
			if s := x.Postman(); s != "" {
				response = map[string]interface{}{
					"name":     "Example",
					"header":   map[string][]interface{}{"Content-Type": {"application/json"}},
					"body":     s,
					"status":   200,
					"language": "json",
				}
			}
		}
	}

	// 结果.
	x := map[string]interface{}{
		"name": o.GetTitle(),
		"request": map[string]interface{}{
			"method": o.GetMethod(),
			"header": make([]interface{}, 0),
			"body":   request,
			"url": map[string]interface{}{
				"raw": fmt.Sprintf(
					`{{protocol}}://%s.{{domain}}%s`,
					o.controller.GetScanner().GetDomainPrefix(),
					o.GetRoute(),
				),
				"host": []string{
					o.controller.GetScanner().GetDomainPrefix(),
					"{{domain}}",
				},
				"path": strings.Split(strings.TrimPrefix(o.GetRoute(), "/"), "/"),
			},
		},
		"response": []interface{}{response},
	}

	if response != nil {

	}

	return x
}

func (o *action) SetCommentBlock(cb CommentBlock) Action {
	// 1. 标题.
	if ti := cb.GetTitle(); ti != "" {
		o.title = ti
	}

	// 2. 描述.
	o.description = cb.Markdown()

	// n. 注解.
	for k, vs := range cb.GetAnnotations() {
		switch k {
		case "ignore":
			{
				if len(vs) > 0 {
					if vs[0] == "" {
						o.ignore = true
					} else if b, be := strconv.ParseBool(vs[0]); be == nil {
						o.ignore = b
					}
				}
			}
		case "request", "input":
			{
				if len(vs) > 0 {
					o.request = vs[0]
					o.controller.GetScanner().AddPayload(vs[0])
				}
			}
		case "response", "output":
			{
				if len(vs) > 0 {
					o.response = vs[0]
					o.controller.GetScanner().AddPayload(vs[0])
				}
			}
		case "sdk":
			{
				if len(vs) > 0 && vs[0] != "" {
					o.sdk = vs[0]
				}
			}
		case "version":
			{
				if len(vs) > 0 {
					o.version = vs[0]
				}
			}
		}
	}
	return o
}

func (o *action) SetSource(path string, line int) Action {
	o.sourcePath = path
	o.sourceLine = line
	return o
}

func (o *action) init() {
	ns := o.name

	// 1. 请求方式.
	r1 := regexp.MustCompile(`^([A-Z][a-z]+)`)
	if m1 := r1.FindStringSubmatch(ns); len(m1) == 2 {
		o.method = strings.ToUpper(m1[1])
		ns = r1.ReplaceAllString(ns, "")
	}

	// 2. 请求地址.
	r2 := regexp.MustCompile(`([A-Z]+)`)
	if o.uri = r2.ReplaceAllStringFunc(ns, func(s string) string {
		m2 := r2.FindStringSubmatch(s)
		return fmt.Sprintf("/%s", strings.ToLower(m2[1]))
	}); o.uri == "" {
		o.uri = "/"
	}
}

// 渲染模板.
func (o *action) render() (text string, err error) {
	// 1. 标题.
	text += fmt.Sprintf(`# %s`, strings.TrimSpace(o.title)) + "\n"
	text += "\n"

	// 2. 参数.
	text += fmt.Sprintf("**路由** : ·%s %s·", o.GetMethod(), o.GetRoute()) + "<br />\n"
	text += fmt.Sprintf("**域名** : ·%s.%s·", o.controller.GetScanner().GetDomainPrefix(), o.controller.GetScanner().GetDomain()) + "<br />\n"
	text += fmt.Sprintf("**版本** : ·%s·", o.version) + "<br />\n"
	text += "\n"

	// 3. 注释.
	if o.description != "" {
		text += o.description + "\n\n"
	}

	// 4. 入参.
	if o.request != "" {
		if x := o.controller.GetScanner().GetPayload(o.request); x != nil {
			if s := x.Markdown(false); s != "" {
				text += fmt.Sprintf("### 入参\n\n%s\n\n", s)
			}
		}
	}

	// 5. 出参.
	if o.response != "" {
		if x := o.controller.GetScanner().GetPayload(o.response); x != nil {
			if s := x.Markdown(true); s != "" {
				text += fmt.Sprintf("### 出参\n\n%s\n\n", s)
			}
		}
	}

	// m. 结尾.
	text += fmt.Sprintf("--------\n\n")
	text += fmt.Sprintf("**入口** : ·%s·.·%s()·<br />\n", o.controller.GetName(), o.name)
	text += fmt.Sprintf("**源码** : ·%s%s: %d·<br />\n", o.controller.GetScanner().GetControllerPath(), o.sourcePath, o.sourceLine)
	text += fmt.Sprintf("**更新** : ·%s·\n", time.Now().Format("2006-01-02 15:04"))

	// n. 完成.
	text = strings.ReplaceAll(text, "·", "`")
	return
}
