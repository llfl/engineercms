package controllers

import (
	"../models"
	"github.com/astaxie/beego"
	"path"
	"strconv"
)

type SearchController struct {
	beego.Controller
}

//显示搜索也
//用模态框弹窗吧_比较麻烦
func (c *SearchController) Get() { //search用的是get方法
	pid := c.Input().Get("productid")
	key := c.Input().Get("keyword")
	if pid != "" {
		c.Data["IsProduct"] = true
	} else {
		c.Data["IsProject"] = true
	}
	c.Data["Pid"] = pid
	c.Data["Key"] = key
	c.TplName = "searchs.tpl"
}

//搜索项目
func (c *SearchController) SearchProject() { //search用的是get方法
	key := c.Input().Get("keyword")
	if key != "" {
		searchs, err := models.SearchProject(key)
		if err != nil {
			beego.Error(err.Error)
		} else {
			c.Data["json"] = searchs
			c.ServeJSON()
		}
	} else {
		c.Data["json"] = "关键字为空！"
		c.ServeJSON()
	}
}

//首页里的搜索所有成果
func (c *SearchController) SearchProduct() { //search用的是get方法
	key := c.Input().Get("keyword")
	if key != "" {
		products, err := models.SearchProduct(key)
		if err != nil {
			beego.Error(err.Error)
		}
		//由product取得proj
		//取目录本身
		// proj, err := models.GetProj(products.ProjectId)
		// if err != nil {
		// 	beego.Error(err)
		// }
		//根据目录id取出项目id，以便得到同步ip
		// array := strings.Split(proj.ParentIdPath, "-")
		// projid, err := strconv.ParseInt(array[0], 10, 64)
		// if err != nil {
		// 	beego.Error(err)
		// }
		//由proj id取得url

		// beego.Info(Url)
		link := make([]ProductLink, 0)
		Attachslice := make([]AttachmentLink, 0)
		Pdfslice := make([]PdfLink, 0)
		Articleslice := make([]ArticleContent, 0)
		for _, w := range products {
			// Url, _, err := GetUrlPath(w.ProjectId)
			// if err != nil {
			// 	beego.Error(err)
			// }
			//取到每个成果的附件（模态框打开）；pdf、文章——新窗口打开
			//循环成果
			//每个成果取到所有附件
			//一个附件则直接打开/下载；2个以上则打开模态框
			Attachments, err := models.GetAttachments(w.Id)
			if err != nil {
				beego.Error(err)
			}
			//对成果进行循环
			//赋予url
			//如果是一个成果，直接给url;如果大于1个，则是数组:这个在前端实现
			// http.ServeFile(ctx.ResponseWriter, ctx.Request, filePath)
			linkarr := make([]ProductLink, 1)
			linkarr[0].Id = w.Id
			linkarr[0].Code = w.Code
			linkarr[0].Title = w.Title
			linkarr[0].Label = w.Label
			linkarr[0].Uid = w.Uid
			linkarr[0].Principal = w.Principal
			linkarr[0].ProjectId = w.ProjectId
			// linkarr[0].Content = w.Content
			linkarr[0].Created = w.Created
			linkarr[0].Updated = w.Updated
			// linkarr[0].Views = w.Views
			for _, v := range Attachments {
				// fileext := path.Ext(v.FileName)
				if path.Ext(v.FileName) != ".pdf" && path.Ext(v.FileName) != ".PDF" {
					attacharr := make([]AttachmentLink, 1)
					attacharr[0].Id = v.Id
					attacharr[0].Title = v.FileName
					// attacharr[0].Link = Url
					Attachslice = append(Attachslice, attacharr...)
				} else if path.Ext(v.FileName) == ".pdf" || path.Ext(v.FileName) == ".PDF" {
					pdfarr := make([]PdfLink, 1)
					pdfarr[0].Id = v.Id
					pdfarr[0].Title = v.FileName
					// pdfarr[0].Link = Url
					Pdfslice = append(Pdfslice, pdfarr...)
				}
			}
			linkarr[0].Pdflink = Pdfslice
			linkarr[0].Attachmentlink = Attachslice
			Attachslice = make([]AttachmentLink, 0) //再把slice置0
			Pdfslice = make([]PdfLink, 0)           //再把slice置0
			// link = append(link, linkarr...)
			//取得文章
			Articles, err := models.GetArticles(w.Id)
			if err != nil {
				beego.Error(err)
			}
			for _, x := range Articles {
				articlearr := make([]ArticleContent, 1)
				articlearr[0].Id = x.Id
				articlearr[0].Content = x.Content
				articlearr[0].Link = "/project/product/article"
				Articleslice = append(Articleslice, articlearr...)
			}
			linkarr[0].Articlecontent = Articleslice
			Articleslice = make([]ArticleContent, 0)
			link = append(link, linkarr...)
		}
		c.Data["json"] = link
		c.ServeJSON()
	} else {
		c.Data["json"] = "关键字为空！"
		c.ServeJSON()
	}
}

// @Title get wx drawings list
// @Description get drawings by page
// @Param keyword query string  true "The keyword of drawings"
// @Param projectid query string  false "The projectid of drawings"
// @Param searchpage query string  true "The page for drawings list"
// @Success 200 {object} models.GetProductsPage
// @Failure 400 Invalid page supplied
// @Failure 404 drawings not found
// @router /searchwxdrawings [get]
//小程序取得所有图纸列表，分页_plus
func (c *SearchController) SearchWxDrawings() {
	// wxsite := beego.AppConfig.String("wxreqeustsite")
	limit := "5"
	limit1, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		beego.Error(err)
	}
	page := c.Input().Get("searchpage")
	page1, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		beego.Error(err)
	}
	var offset int64
	if page1 <= 1 {
		offset = 0
	} else {
		offset = (page1 - 1) * limit1
	}

	pid := c.Input().Get("projectid")
	// beego.Info(pid)
	var pidNum int64
	if pid != "" {
		pidNum, err = strconv.ParseInt(pid, 10, 64)
		if err != nil {
			beego.Error(err)
		}
	}
	// beego.Info(pidNum)

	key := c.Input().Get("keyword")
	var products []*models.Product
	if key != "" {
		if pidNum == 0 { //搜索所有成果
			products, err = models.SearchProductPage(limit1, offset, key)
			if err != nil {
				beego.Error(err.Error)
			}
		} else {
			products, err = models.SearchProjProductPage(pidNum, limit1, offset, key)
			if err != nil {
				beego.Error(err.Error)
			}
			// beego.Info(pidNum)
		}
		Pdfslice := make([]PdfLink, 0)
		for _, w := range products {
			//取到每个成果的附件（模态框打开）；pdf、文章——新窗口打开
			//循环成果
			//每个成果取到所有附件
			//一个附件则直接打开/下载；2个以上则打开模态框
			Attachments, err := models.GetAttachments(w.Id)
			if err != nil {
				beego.Error(err)
			}
			//对成果进行循环
			//赋予url
			for _, v := range Attachments {
				if path.Ext(v.FileName) == ".pdf" || path.Ext(v.FileName) == ".PDF" {
					pdfarr := make([]PdfLink, 1)
					pdfarr[0].Id = v.Id
					pdfarr[0].Title = v.FileName
					if pidNum == 25002 { //图
						pdfarr[0].Link = "data:image/svg+xml,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20width%3D%2232%22%20height%3D%2232%22%3E%3Crect%20fill%3D%22%233F51B5%22%20x%3D%220%22%20y%3D%220%22%20width%3D%22100%25%22%20height%3D%22100%25%22%3E%3C%2Frect%3E%3Ctext%20fill%3D%22%23FFF%22%20x%3D%2250%25%22%20y%3D%2250%25%22%20text-anchor%3D%22middle%22%20font-size%3D%2216%22%20font-family%3D%22Verdana%2C%20Geneva%2C%20sans-serif%22%20alignment-baseline%3D%22middle%22%3E%E5%9B%BE%3C%2Ftext%3E%3C%2Fsvg%3E" //wxsite + "/static/img/go.jpg" //当做微信里的src来用
						pdfarr[0].ActIndex = "drawing"
					} else { //纪
						pdfarr[0].Link = "data:image/svg+xml,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20width%3D%2232%22%20height%3D%2232%22%3E%3Crect%20fill%3D%22%23009688%22%20x%3D%220%22%20y%3D%220%22%20width%3D%22100%25%22%20height%3D%22100%25%22%3E%3C%2Frect%3E%3Ctext%20fill%3D%22%23FFF%22%20x%3D%2250%25%22%20y%3D%2250%25%22%20text-anchor%3D%22middle%22%20font-size%3D%2216%22%20font-family%3D%22Verdana%2C%20Geneva%2C%20sans-serif%22%20alignment-baseline%3D%22middle%22%3E%E7%BA%AA%3C%2Ftext%3E%3C%2Fsvg%3E" //wxsite + "/static/img/go.jpg" //当做微信里的src来用
						pdfarr[0].ActIndex = "other"
					}
					pdfarr[0].Created = v.Created
					// timeformatdate, _ := time.Parse(datetime, thisdate)
					// const lll = "2006-01-02 15:04"
					pdfarr[0].Updated = v.Updated //.Format(lll)
					Pdfslice = append(Pdfslice, pdfarr...)
				}
			}
		}
		c.Data["json"] = map[string]interface{}{"info": "SUCCESS", "searchers": Pdfslice}
		c.ServeJSON()
	} else {
		c.Data["json"] = map[string]interface{}{"info": "关键字为空"}
		c.ServeJSON()
	}
	// var user models.User
	// //取出用户openid
	// JSCODE := c.Input().Get("code")
	// if JSCODE != "" {
	// 	APPID := beego.AppConfig.String("wxAPPID2")
	// 	SECRET := beego.AppConfig.String("wxSECRET2")
	// 	app_version := c.Input().Get("app_version")
	// 	if app_version == "3" {
	// 		APPID = beego.AppConfig.String("wxAPPID3")
	// 		SECRET = beego.AppConfig.String("wxSECRET3")
	// 	}
	// 	// APPID := "wx7f77b90a1a891d93"
	// 	// SECRET := "f58ca4f28cbb52ccd805d66118060449"
	// 	requestUrl := "https://api.weixin.qq.com/sns/jscode2session?appid=" + APPID + "&secret=" + SECRET + "&js_code=" + JSCODE + "&grant_type=authorization_code"
	// 	resp, err := http.Get(requestUrl)
	// 	if err != nil {
	// 		beego.Error(err)
	// 		return
	// 	}
	// 	defer resp.Body.Close()
	// 	if resp.StatusCode != 200 {
	// 		beego.Error(err)
	// 	}
	// 	var data map[string]interface{}
	// 	err = json.NewDecoder(resp.Body).Decode(&data)
	// 	if err != nil {
	// 		beego.Error(err)
	// 	}
	// 	var openID string
	// 	if _, ok := data["session_key"]; !ok {
	// 		errcode := data["errcode"]
	// 		errmsg := data["errmsg"].(string)
	// 		c.Data["json"] = map[string]interface{}{"errNo": errcode, "msg": errmsg, "data": "session_key 不存在"}
	// 	} else {
	// 		openID = data["openid"].(string)
	// 		user, err = models.GetUserByOpenID(openID)
	// 		if err != nil {
	// 			beego.Error(err)
	// 		}
	// 	}
	// }
	// var userid int64
	// if user.Nickname != "" {
	// 	userid = user.Id
	// } else {
	// 	userid = 0
	// }
}

//在某个项目里搜索成果：全文搜索，article全文，编号，名称，关键字，作者……
func (c *SearchController) SearchProjProducts() {
	// limit := "15"
	limit := c.Input().Get("limit")
	// beego.Info(limit)
	limit1, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		beego.Error(err)
	}
	page := c.Input().Get("pageNo")
	// beego.Info(page)
	page1, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		beego.Error(err)
	}
	var offset int64
	if page1 <= 1 {
		offset = 0
	} else {
		offset = (page1 - 1) * limit1
	}

	pid := c.Input().Get("productid")
	var pidNum int64
	// var err error
	if pid != "" {
		pidNum, err = strconv.ParseInt(pid, 10, 64)
		if err != nil {
			beego.Error(err)
		}
	}
	key := c.Input().Get("keyword")
	searchText := c.Input().Get("searchText")
	// if searchText != "" {
	count, products, err := models.SearchProjProduct(pidNum, limit1, offset, key, searchText) //这里要将侧栏所有id进行循环
	if err != nil {
		beego.Error(err.Error)
	} else {
		link := make([]ProductLink, 0)
		Attachslice := make([]AttachmentLink, 0)
		Pdfslice := make([]PdfLink, 0)
		Articleslice := make([]ArticleContent, 0)
		for _, w := range products {
			Url, _, err := GetUrlPath(w.ProjectId)
			if err != nil {
				beego.Error(err)
			}
			//取到每个成果的附件（模态框打开）；pdf、文章——新窗口打开
			//循环成果
			//每个成果取到所有附件
			//一个附件则直接打开/下载；2个以上则打开模态框
			Attachments, err := models.GetAttachments(w.Id)
			if err != nil {
				beego.Error(err)
			}
			//对成果进行循环
			//赋予url
			//如果是一个成果，直接给url;如果大于1个，则是数组:这个在前端实现
			// http.ServeFile(ctx.ResponseWriter, ctx.Request, filePath)
			linkarr := make([]ProductLink, 1)
			linkarr[0].Id = w.Id
			linkarr[0].Code = w.Code
			linkarr[0].Title = w.Title
			linkarr[0].Label = w.Label
			linkarr[0].Uid = w.Uid
			linkarr[0].Principal = w.Principal
			linkarr[0].ProjectId = w.ProjectId
			// linkarr[0].Content = w.Content
			linkarr[0].Created = w.Created
			linkarr[0].Updated = w.Updated
			// linkarr[0].Views = w.Views
			for _, v := range Attachments {
				// fileext := path.Ext(v.FileName)
				if path.Ext(v.FileName) != ".pdf" && path.Ext(v.FileName) != ".PDF" {
					attacharr := make([]AttachmentLink, 1)
					attacharr[0].Id = v.Id
					attacharr[0].Title = v.FileName
					attacharr[0].Link = Url
					Attachslice = append(Attachslice, attacharr...)
				} else if path.Ext(v.FileName) == ".pdf" || path.Ext(v.FileName) == ".PDF" {
					pdfarr := make([]PdfLink, 1)
					pdfarr[0].Id = v.Id
					pdfarr[0].Title = v.FileName
					pdfarr[0].Link = Url
					Pdfslice = append(Pdfslice, pdfarr...)
				}
			}
			linkarr[0].Pdflink = Pdfslice
			linkarr[0].Attachmentlink = Attachslice
			Attachslice = make([]AttachmentLink, 0) //再把slice置0
			Pdfslice = make([]PdfLink, 0)           //再把slice置0
			// link = append(link, linkarr...)
			//取得文章
			Articles, err := models.GetArticles(w.Id)
			if err != nil {
				beego.Error(err)
			}
			for _, x := range Articles {
				articlearr := make([]ArticleContent, 1)
				articlearr[0].Id = x.Id
				articlearr[0].Content = x.Content
				articlearr[0].Link = "/project/product/article"
				Articleslice = append(Articleslice, articlearr...)
			}
			linkarr[0].Articlecontent = Articleslice
			Articleslice = make([]ArticleContent, 0)
			link = append(link, linkarr...)
		}

		table := prodTableserver{link, page1, count}

		c.Data["json"] = table
		c.ServeJSON()
	}
	// }
}

//未修改
func (c *SearchController) SearchProjects() {
	limit := "15"
	limit1, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		beego.Error(err)
	}
	page := c.Input().Get("searchpage")
	page1, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		beego.Error(err)
	}
	var offset int64
	if page1 <= 1 {
		offset = 0
	} else {
		offset = (page1 - 1) * limit1
	}
	pid := c.Input().Get("productid")
	pidNum, err := strconv.ParseInt(pid, 10, 64)
	if err != nil {
		beego.Error(err)
	}
	key := c.Input().Get("keyword")
	if key != "" {
		_, products, err := models.SearchProjProduct(pidNum, limit1, offset, key, "") //这里要将侧栏所有id进行循环
		if err != nil {
			beego.Error(err.Error)
		} else {
			c.Data["json"] = products
			c.ServeJSON()
		}
	}
}

//搜索wiki
func (c *SearchController) SearchWiki() { //search用的是get方法
	tid := c.Input().Get("wikiname")
	if tid != "" {
		c.Data["IsWiki"] = true
		// c.Data["IsSearch"] = true
		c.Data["IsLogin"] = checkAccount(c.Ctx)
		c.TplName = "searchwiki.tpl"
		Searchs, err := models.SearchWikis(tid, false)
		if err != nil {
			beego.Error(err.Error)
		} else {
			c.Data["Searchs"] = Searchs
		}
	} else {
		c.Data["IsWiki"] = true
		// c.Data["IsSearch"] = true
		c.Data["IsLogin"] = checkAccount(c.Ctx)
		c.TplName = "searchwiki.tpl"
	}
}
