package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	_ "github.com/astaxie/beego/config"
)

// 課程資訊
type iClass struct {
	ID       string
	Name     string
	Location string
	Detail   iDetail
}

// 課程詳細資訊
type iDetail struct {
	SignUpStatus    string
	SignUpTime      string
	AttendClassTime string
	StudyHours      string
	Location        string
	MapDetail       gMap
}

// google map 串接資訊
type gMap struct {
	Destination_addresses []string `json:"destination_addresses"`
	Origin_addresses      []string `json:"origin_addresses"`
	Rows                  []Row    `json:"rows"`
	Status                string   `json:"status"`
}

type Row struct {
	Elements []Elements `json:"elements"`
}

type Elements struct {
	Distance Distance `json:"distance"`
	Duration Duration `json:"duration"`
	Fare     Fare     `json:"fare"`
	Status   string   `json:"status"`
}

type Distance struct {
	Text  string `json:"text"`
	Value int    `json:"value"`
}

type Duration struct {
	Text  string `json:"text"`
	Value int    `json:"value"`
}

type Fare struct {
	Currency string `json:"currency"`
	Text     string `json:"text"`
	Value    int    `json:"value"`
}

// google map 串接資訊 end

type MainController struct {
	beego.Controller
}

func (c *MainController) URLMapping() {
	c.Mapping("Post", c.Post)
}

// index 頁面
func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

// post 取得課程資訊
func (c *MainController) Post() {
	var result []iClass

	result = postInservicePage(0)

	c.Data["json"] = &result
	c.ServeJSON()
}

// 爬蟲 取得課程資訊
func postInservicePage(page int) []iClass {
	postData := make(map[string]string)
	postData["__EVENTVALIDATION"] = "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"
	postData["__VIEWSTATE"] = "XNurmubdqRP0G1j4eg3LgJLfiA+jO1KqLe9ehO7mKZKvij4t9uCAOcf61+/pLReUikkm7V9TSsGm6lvWepsCAD/ECdflQ4IZuLUSDGF64PTiYS974hptSq+ES8I0lra888OgUpGSH73Ew02S8ZNCbslodn3Upx4aCQnK6hn1YMO0hjNb9vOY3m6Ey3QV1/qU+jnUPlYInLi46892FQ6ttQWHdy+6DizRBPIom3tdotiR+MKi0ubjFHWs8WsUnrNt+M/BxCswSnI4GuLJ0KD02lmJY1+VQ4QRXwGNNa9MPhURbzHTXEOGCjzyVxG9A1w9CqrRxxoV2lJfdLhwbhMnhZapwFPVNCWPb4KuAm2s6dfl0pK5u5GXRHN6yaVTprCIzUXtcM1RVDX+MSntl5k7j64CsTx3fjITpVcfYtHDvJE3TZNmw/M5mRliVveXcXYjuocGW7QnVV5BUwctrMD5iyCPRA/9EnYh/ZcGj2zwxMRb857UKFalgwnh4xrfQKIW+i6j/HbaPC67xvGlazYBCAvbzTGOOr8E63oXRBflxRchILNoZHo1YCzO8NQT6BmTKsnr64plPPHPutBviMLpI4k/xbZXMhdexZ2VAjwDVT0T+dYBULabKzt+iYlm+2BSx0U5q4x8IS2BK6zlNH/W0AJ4bJ5vEs1D6Jbrhl4dDSFrMD78hv/R68i0suMfrgY8ZKNAYhCBFTg1bAld8/zGdCJEkyulq/T1D/4ZdUsveD7XsWEgbxThqttu+w1xuWlPKhT2FteweG1lhmkm/0Agnd+1Xtcg0zZ5qTBKMpTFMut2dcufRdLNi5ff0inK4c6RojIdQf6FTvPpSsz7ySuPFN+LJGJZZVpfx1CiaWStwf9Y21x73xK34kdHebKEAO6gfM0euOs3MZk0YNB21LetDPie0pq1uIrhh+N19oTpvIRKeT5B1JMc6ekXBlGJcj1rIBrmxO9L33pYJJdM2heLFCfTdUi8j3G5ah1yetNbe1zyiRMBtMqf9w8OYbqCNSgsTwWlFDQRvHwEFPXofgoIlJNCbkkB19OSKBoyX74725klzEmBrKm7kpoQzHf+bu/qoFRU/cBqaiWT+VwUuUN7T0JzDyYT+K+I/ZnFNDOQrDzQY0e2Rzj0Qn4F8sQnMVzePHKItdYjeU0RBfjDbysoSY0KIC8fAOOd6Hq59GKJP/Z1BOEuPf/c9kYZtqrvkSSpu0Ux3Xjy5gdlM3kNfVbkxAKN7xSi4VOJOddVM18T8Zy9aJfUfW6VvQksFGKPsOlhCmQ6vrkdkOX7uU2LzfhGMt6e0FYNeNizutjxo8nVFML1055BDiVUrhiIk9vCivyeRq5e0hADv75z2k6Q9oIaX1x7C5cnkeTWEVbnsWZepDC8XOg8we0ylixz86UWDxaD8Ceb2ZIYGDVRcx1zhCYfEE4PaDD4An+VX046GSjRsAzgNzaMo0a9XMJQDYeyXLlx7uLLu3qEmSFWI68XcG+X9lTwDsIwwObdQ8+/XOcLswQ0RglpsgsPwBT6GiOUiiWtU0IgdjwYr4rApRtn/hNJt+1j/LnTxgKMrqtoCVSg6NkMh82iVAJ6tVwrqvSJh0UoOZCZqwQF85RCPqr9TZAWxVKmla1flL+aiWbuqwvrBbQljmUTMi5Gp6verDJqTR5NmKS5XXHAZYbiXqx7N0vxH72o360uqzb1/S5zmIYewZap+59bolu48wUl+xSEmbzTJgWW2LQlCHlhP+21gUughT/2Nz56GE+d8BfCNcvzHaP2Uw0mBENThnNmJUIxm6HsiZLQ/EvCn5fSEokWjMPky8ve8eIOK1VnZKsHSqiU6ATshSwsRBTdAADbN4mTgHf+D7EEO9/Ae8drxlneHy7/y66qfFyu8Bb2v2SG9qhhMGngsGF1gbfyWwWK2cO6WDr6Vor8UJyvbOICGGrjeg98ulpVQkOixkYaSBMrG7iVzzBPcQUKJjjaVgWrg0SXUccUIVNhItViDb8F7KpnZOaJiKCi6d6pogcmGRo1N2KyKVdh8wjK85AxVBPWvXWedRCPt2LDQtzaC6FCk0lzUIRWjAjE6fatO1SBvTJqJsG6c1gsevWCxktSQHjgQimCVw7AOohYAQrlSZFsxXgAl/WG3dLQcGGer7sIXHDsABnfjijZI+WNruyy/MUMFmM3JmQZ7fb0uX3/MU2Gc3coO+fF+5sx5kXRS/sDPJeXb+4DFiR8H8LCWc8795jNak93A9qEDZ0d4lJUmjbKI0iyV8sbF5jiY50IK30PsXO7w8cjsCXzi7hUtHPbOBTOmkXEWlWMljeaRYNNRhx6TnRXdcn2ANWUXBR/7v+nt8l5kGpaQl0AuU41YjydeJzXM1AVxR81CwplPw1zad8sqGFvtS5DwLDtZQagEvSS9iNq8JFHyv+W7QHyIQ3U2qImh4Q8d8h7tfqz/6UuFq7bm+Xr9d/vWf08RmM3GCr800q/oek7rs1HwAPjZaC1N60JSr2MvqhGIAiz9mJdNaLIqEJBimE5uhggIs3X2DL5cFlwq2PCHkQKtmIaNUjq+8RizE0g6inSj8jb2J8kdp8cTuj2sAsQM8ED2P01b8zbfzOZvPPvte3Esq4JZIPDtp5ipx7QqLanRDS9jDngmuHwMx7JgSUwRPcAZaKNlKdGno2OFBI38qfa8YpBwhLp0jCTosWRubbTAmizvMClq/Tp3jkxjtGBSpJxMI4nmIQsgFNKojQ9zhUjNNXYm3LI5f67D1cVq+DnqveMFyote7qPwWlRFO7A1e1rVHCj0kVV+R2lfYjoz4MxzNXczhozSKhbCdW1ktWnKYqs9fBwtZxgkq48xuRBUXZN0A8h7T6x17lZpPA4lPdEtCgK+1aAlBBW7u2Kw3Dq84EXypYqk4kLWWMaPOb7B/9ty8I8C29zr6agCvpRRm6+yYTyg0tMr0yEc2opsxHoCL40WnsbF/wDJnLqorarMbxzfeKlxoYODKWEgC7xkFn8mKZJovB4kCKYTX1jjMZ8iePR36bSe0CZrLFXia/n0X384AjIaQiXsskOjE0rnfRkLNRdckEEX+96LzRQOjQFEzHeUP0Z7yHCAIaFJxBsWS5/u4ouieTN+uXv35Et5n098L1pSKkYCzhZP38e1rNduIZYNpseIvlvXd8mVcVQRFf5cuXlnZfnUaXhnCwmb2OnRuA2WD8js4h78qu0ddrheNLFJm6DuTFIZ09OhbBgLN4AIM5pUfeCI71XsfehLtCKlAvhcKwLoTdok6F9J71cnUVWAh2RD6on1Qp1kPTqT9kJ0BYtmPRB8Eamc6XKDjLqJ32KWoTfrxkFmzIuQN9e/q4TDZLIhHW/aRW+D3PrZaL9vOm0aWmaIZqof9uh8xdCd0c2JjAyPRByLuhPgNvIeFEaqx2ajP2f7PMfXjRnqS187IBmLH0QqDhfTWvepQigsAV/0KdwxShb6aSsnxY9OmK3k4iwDi+jrtS4ruzI4rYwXAUJnU+KqjkGWXmRCjNxTSeIdYtE51OFsR8SA4xTRLzA8qew2/5y9L4ql0ohHCfnYwlO5wQQFl31RZgOplihwkxiXpA4Bn0dBnb2GdD2uQ/bAwTYMpHLe07jz8MV4KcTW6KHxD8nzfRD9EgiGksG9waiTjFSC/IcDF4xUGCaERXABwQi7x1iPk9sHX51/lWTb7kNMzDzfdZR8DamFs1w+phlI0AlUDbr8bqe4ezdUsMHv6xGjHJzVcRvGaz4JKLtajds1H70EmK4pjOLkznAG9Epskr8GruKZVpYuCv5Tm8BIek4fo4m6xnIs7G93AHtxf9hMA1GSopHS93EVZXoNU8ZxUk6s2a/QWwWdQaNDIoNIFswwIU6/+dYWsgaWmk6nuiIfndRq87/TPAEG4GWuwK8XRylaXVy/PfPqrOl2Cbtxwvr1wqHZOtT7jhQBDUZsIKUQh1QEdRT3lcrDGIsaFV7Sub9p+TitHzvxTy7PjGpg6Zb8Gs3N5E6gg1rxjdQrjLz8Pojx0FXlYvzrqvWX865hJEmSA3AhCJ2S6tQo3c2wjQHxL/6ojw9rZgKLYGbL1ys4RaAxahAczRBlY+bwkTH7aP3kLH1VGeR27G12L9FekqIrfVQGRDggAbOZuV1nsPhuyzwsLhnEB7NwdrzdacpZ0jTS9zO/ndPErIfK6yH2kJ6tsbX++7AieoG9egOk+Q75H2EVAen5J+Sn0ATMFOjeECaBdgZWaFDFt643Dtyr17fzmtNv4O3LychsBmzbFSXAwu/ibL7v8P7mxca8wNvVclCfRmc+Zy+Xcx9IhCN4cndYYIB9Zsrm3oJ9JtpvfjWzR3CSoagskogIwdPv1x/OvKzPbY/Y23Nez/oq1JR3Ymt8o0rcfs8tX0uvav/CZid7/RM1HgweNsQa4xHUDDt3OohFcLDJTjWWm9++Gzy5DppEEh5Mnd9EdmE+x3JPtqxpywqmFvdCxM0gqc2+vCjBuZtdngpGONiwK7S5Kt6VBV9ebM/oBxmtUohEh70wZp+C/mIgo5KmwfMhjDFk6QMaiWauMR4XnydKD3053DyU9TMWsulkd398jwYvL0wCv7iDEzchMAqR1yzsHo8l+UD1lPJQSYZFmJBRQPBEBZX/C9zquwx5GBNklvCSK+Sexjhqt+w04nRDTkdX1y+AmMwqusG5BfIxgxm9d+mVHdtzsqbOcmG6TTiKvAgjIYWxy2Rj/LLPfXujKQheO1xa+ZvEWtsg5ZHt3hAdvSxF5g7Li7bJsjVTHCQc7tL9ZJqP3zB88qePe05MV2GRsB0wnavxdCz/HDAi/rXOIUYC91yiz38FcPJuhLdGiGHU11dtBD7PQ+pEYagt1GSLvmjX8DILyO3GY5N0DAFw4EdydHNbET6faFHnZZVYlBKeI1K3KB+BlsXGl2YydP2mcXNCqY+SB+UHTD7p79PJEi+FdvLPP28LaJAQc81fB5Bel8KmeHIbDQ/Ahc6tKAdGwbFsVHaHC9Zh0fZOsKyJ+Bc+UlaM6Ja6vQf4bfA95gAkkTNnoz1WAT2onXmLGevJSFE0tX7wBLpopDK85AddXsHCrzYhqJ+3s89clk7yofthA0JQs2gMHIt2OfsWINxBxBPMw0WRO6B7Lb7uM8hC06RYbgjiotoueCVpeLVJWBw6QHFsxT2iqygI5MxpilMrOgAuJfI21ccQV9JPnjInUmVZqHUJsJ7j5aBpdtyuwhmhlkUsgypY1B2cxNokeAzoNunj2RNQ0KdqAQCSwu8qTzRLgkGkRyYZ2FVTJd7WmZ3+zOF0qr/Qg/NoDFBkB5NpNg83vUBZAQA2iQd9dE4aTUTviMjqNA4GZIvhHAiQQh7BzFMhVXUqq10MSuXL4zi9tzPQeBJRwg88AU7zG+FaAdf94dJJCrGsTLzXEPDOvmnPJSLsWTZhGI11Zgm/j92P2mDiLccGbw7D61+Hv29dATkpYEFS8RtFkZxnuf6oy+dFTg+lEq4Ij4Hh08EHF1HtzqTh+5yoTXABuqCcqS/WEjj3L6lNGhyRvPoek90EPA/WLGVR0+FVQpVXCiv8puzcBpAK3icqgCIwDql65re3SGoT1bA+HIrMg10CES/7QQxcq2xZYLgdFZ20i0fMCbqDjbZzHEFl9r1kNrT3Va4Qz1nMObWexu/kQBt+YqDq0Oy6d4/ouLv63cGd5M/aYlSbuDbK+GNNj0T/oY/ZHxSgFeNLl5ym+/w0tSyy2fUSB/yJ+AipMY2epJZA61pNGqjS372FgOlmS2sWiKsc1EknVTe9Ymx/Fy04JW1TphXzv8kWZtYyp2saM+ihtQ/zriyfGvU4n8EkSy+KoB6+AZK47WtWSjCj3/GndnOJchFHhbaJBRod/bbav7Td8RJwRUKzohls4uFqWDl65Z6DBlqtJuIR0YkYGbgo+R90CJmRqWCDsnZVGHsWXAhwY5JJn9QoC9YGzwSRMK7gqCi5ZGCcc7KpkN8HYMK3dvSM184c4QZfjH5GpRRrBIMS8/oYQPtkHDZtWHAr4PVCf4GfmRIZKa+vzetqB+2W72t8eEWvsnANlqNnfna0Vg8pM2qZYMU/Z69d6Sgia7Jp3gAT5ICyjebwYe3c62CSJwzDxhdQXBDqvYUW8/UR4+YOofgrYmObaLOXI9tISvDDJJy+HhLqBgeiM3xydJkcV8QPEhtnczy9zWsoxpu1JcOaX7v5k077M6cHS//0+0tqlyyt3a1KuHuTwTRSjUa9Nh55i1yk939eUVeXVOMSlrKh81GFBga2i9UB69YnxjrJZfonOG/Tn5Eg/+yVuJmfliTt66/0xK2olU2irTcobivWto0OY0QkBGYtEFsTQSmq//vO9QZB9IDq+d+gMFswrqSF5KT5E4aoCDxuzM7hc507pRklle1FIPpeWfTGWD8tET1KZZcymZ7cJdx24KRRWCkMSWA2ww9HZZFBwkEVaEip1LnkKUZyX3J+1OyB0wvy2f4juGMjXO0yuLwiyXzhos/7Pvc4MyRPQrMDOD6LR+O4hGi0KTYwv+qIKwRWRFGvAt97wS06BfKiKDF+A8EcIhPrx/rz6kJMm8PcFNfopAWH77HpdyaEZ5uAG8v4lHh2GTnxpwZ9u4AQIdLSdnlvODcyXZTBfU0p9L4M2MEr529+5F1wWe/nN56J4D5O+ln78VQb0gfdKNsbWk+G0iWYmWCcMfHjaJDEpeNq3nR9PQGIIkLqv6ZDN2efuC+Yb8XGLahPyprXeO43iwUDTff+THWAD9PAzRaX3w9KHeJca4crBCO49WFbcEu6oKBAxt2rvrOKufCtA79UBP6L5xvAK8LtHpL7JYrM1tQB6TICh/OCl8BI6I42vT9k6Smo8HYROJbbIh1Rnxui0+PW0Gf3v7WZeJnBJobNGX5gsGJZxSN3BUwzvj4J2Ci/Ayb5jsdXAmjG/q2JVQM5lBW2stwQSDA/k2kDpDiZztXx+TgL9SfWqA2WsHLV2aLVq+AmCe/gg3l6wD6/612naNmaviQTJqJtfKg7r42zDFg38BohnTdZHI5sgyT/ay4dkx8y/5lwRnti4awuh/tiEtlzeeKrpqZ6cftR0xoP75ZF5UdxpLsd6K7rQy2rE0yuzLs5nAPgaoPo7agXGwVoFYAUEGABmB3OqXb0w08IcbvyMSwZSxkPf2VnYTWkDIbpuqYPh9aC8Owg6/a6LuXy6apVQONlgCwvAFLwAtj3ef3+ZimKwG5GVz4Qo2lZwH38yqkU0uoj4+NwPNXvqCFuJf+9DM+L2llbyb6cTCRH1DG3I0nJ7HdB1b05DHIKUjFdld/Bb3Dj9Zc+xJS4KmnF0sNBnUX6AD3TUQrMHEqoHwg84V+eWmbAmWhhtOuTDSYdY9J4pORb7cKMsw3zM0somgk9jCuUMESzqlG1vuL+zDPVjSwjxMQY+h6NBvumi7ohts1nFlDC2kTnEUTlaBlm1A+69LSi//fLaamrv5RSviIxItzxOvPuatRS2+IJljE6uHVttEMOWxeRNnh19O+5yEKNFVnHg7JuYXlSJx4xMglqzBBJ17H7mmaIdSZ9G7h3novi391NUz/opUJOn+YjGUYTr5byG/9hr3AokevrDKHEsY97coeYbe4iRuJO2OT+kxHSmVKsQI7EhkTj8P7XF3peAGzXcq0mgpZTZ3Oa05jVSxWOK2cdBQXepa4990IuAEmU7NJJTaNR/Cdr0yDX2jWhZIHPH/dA0XAmrXppltSourwzOSHrRUtdPT9/zHdrzOV5QoSuA1f9+l8Ig2qAm3pNw9sG+w0AZfH54PU3k1Cl2RE3PQ52uqPUghmf9swLB/wkaQsXIhKT4mF4KDbEdTwXcZ80mlMRgiSC5TA0+wsm2HV0HiP2re+AHPVyrEM81lQ8FfzH3Zenf4EhAaj3GRk5lHgHWqYm5uzafWRvtoLMo5ylh7m7wB90zu4h6FNLi0lOSySHO5lIctVCTEfI/vpRZIgvifCY++ieFOTyQ7LDubwGNN0CdJQLsE52lPMdR+UlqDgYl13eY35aJjrWVIaIPHk34sm/ue9F/MX1rfRwFQJ6eTO0IGsZJt9ZOEjvmTjNJlXSHbTOAR0ICjVwPWVWBijcYLOPply/C+5HE2E9AobhMEDQ99OGASantcZpRnoJ0Cywmyn82QN71LcyWERXMizqEeHyTKhJKfnPhksO6eduAvbrgP0FtyHio5S2DdsIO5uHfO8UNc6m42odKVXgTvycTs9UdoBa7K7FiZqT6wJquS0N9XnFEareEw8F3pRNZDChp/n61CMYvOr4vpTtNqxG5sqioUy2tZI9dOULAnnVonZwiMDK3hBXHui4TSaltlX48pwvzjJ+49w7i6sAVBUeDsYs4vdOuXyx820wToMGaaclgVE92yKE74XS2F97lutPyp6Y7EwR9FuRHZavQbLX9X6MOesspP9thG4eGMDCIFXVgpOU7AotHGsdQqanIIiK5Ds/v7/cQp3TrmjByZormCz4D+A6GkU0ASKxkLcv282vIyM1Qbv6DVaXX3e4L3M9arfuN4SjUECD8QCvIR8yR+TNd/AnydXk0kCECFx4u2zcedtJj0K3oA1YZtoDUya2ud6B+229F7F1id9n6GgxcpqRCxXTdW4CY5v/RGtZDndMuvDZjuEHPDZ4oi+rdey7H5ce0rJy89GRDQcZ77TOUueOWeeK6T3ZNIPOzPhytb8AnQYmc8qjHGTwOPajx1BolG/hIF+KjQcqGubqoypVKDY8FDWt+U9s0sikxYr1LfS83uqcM8ldAsuaB+lucZGVZinElVmR8KcQ4rQ4w2DWaCABCm6d3otkYl1UVdn58sERCHanoApK9U7Ve1BzgXL89XDJK2W+lxMil/GWX5H8+3F3rDWlgjMNXVDWHgnDx3lT8d52xV6iGYHEE5NbWk5beZXp2yLMT5YY0S84IXNpm/YMgvVt+sgrvZthYQrAtO3/h5h7Vr9gJbPtqsh2cU1zwi6qX4P+SdP3HFQC1y4WvziGwLGDpAdZk0OlRXzRptE6S8fsnlJS3WjusK/OQ6MfbN4K0v+5VlamxYmt2iYCYC6k9l6YEntSQQlZQgUOkBKSHCZzLr8st7mvYcHr5uQdriECvwImINcrBR3Z6JTKHBCYl0KStPOQUMrBYtX0SSrB/QhowLu3tklNFESEp4nqAGyd9fGIQlLSi9PJtFeMMM4EkHbj7R+V7QrmfcwShIYEvCqwhUlcGYviUHuT50aqXar7BqJTWhYxMQTcX8NVgU+acesff1m9AiE025ACY/nVJsKIxYaPJv+1roqtlJpzs8rBOkwZReODmOJ2fMjrUKPTypWxisodc7pXQYdUe+of8UIYZA64dxlN2d4VAWWUex0bsDbof9Zp0ZAqYLNFoiMxICUC7UBQGZ9D+EyiQf/769j/XmKs//n1jXaLs35OS+RThkQuqzXM1/g/Q9Kxcsr5jZegRmyQlXGbq/n1FG5+DSrgRZgsEhJnWIj2HbDLgqSXGG18CKkhBP3VYcIoe2VtreOTqGwJL11BT1MxTBqDWjQ+VSoScr5fqMMvMKGnSIt5Y7hXoYyk9CxvJ4NtLNzuygm+r8bea0YU3ke2Zsf1SXeZHlM6GMKSP2Hi3YBht08YORWRLsTjuR8op+vnT0UJeRg3WxfAfv5LqspBg4SFGxTEead+0JtvkcD7kRtpH02YMH9Vi24hznlwttbfbslzuQVVQMP5WPBzdoID6ctp0fu4u49tgt3dECS+pRIQ0mgGRubhWbzwNuofx/lDnhHDv7Vz5jIX4eDx0j9l5z61HPBEP9FVd1Gm/LqD2S9murygoDJjyy1qq51L3+F5OABQAymZMzrQs2lvk5ljqUpSm3+1QFoUvfMSarlYqejP3Yq20XdpiTN0POFBZdYEPwpiBhwPjLLAAfmlF6g2k3WqwNNAcjoFbw85OmuaoopjQI6rbf7PacUPBQ/hOosBpswIbYvb2YbgnADDA0EhC9UNkaVtwQfqg8xfiTIlUKs5IroVvxJOkNP5AE+R96L8drCbgby7QxRBB59kFT/XSORA804IJYkreJWnDV8xokrHaEo+ousyZcyfoMUs0FDj9VksudCfXVwwZFmqBvztsMBO8i+EVVjoejyyJQaDu101kHJIAsYxu0Q2cTSI5vb8v557v555gd3IgBIiFeyXh6u2k/xFbGpemfZ5ydHRu5r7R7l+5HbhxSOLVlcmi2RE53POlWy76hs6Xs6pvs2NeVT3Y0Ppql7t22wi0lJSm7HjY0WP0eJXWe5hVykLBsv/gyWZ+pd+UDAuzLUEkoXsgdlDPmgsDB1rvKN2vWzuNkOeKygO9ymYOQQhtPmHVtErPPBC7h0QlRO4U8xjRgjFPPfVVib/uvEYjYTL+LrU2uGkqBpZxbVe/c1BQ47hnyG0zc30wauW/FH4osuRW/VQvNXTs23RHGhN/VR5HbYfnhD3DTFQsWYfEM0hALeWKat1eDDooaA+rYnLy9QNMVzKKPxA3oYXSn2pN4FK9oBjJoAPb1sOxqFj4a9mZiZXz1L679Ld3FyXTL18n7/avlvsFSqwk64oVv81RbRfXpO3GbG7HJYZl9BwkG5PKkjByB/ztUt27T0UHH2Ho3Lha0rj48CjBtl7JNYf6bVSapDWv/7Aoi+rllXiEqC56gmrMyYshbsl3lIjL2VPtJh2+cnM1lZ9mVUZWwXvFbBqH7Xao3D6qiNVRgrBFIHfY2JnSBu0SWHFtlLCUn/pJBDjY3LY447vMxZu6Icz5OJ2nP1H7VgaeGAY4tXcbs24yIHWNP8OtTGyAtQ8HSZ9vz0DOdAcesaBkCaDsZN0FwqB4NCjfBTCMbN27xvCYiIpnc2VKQe0Ff80iWylc5QKuKJmA9YA6LxMC1ZInmEb0xFmeYvXAQxt/H1fHHuljQq175sx0ufxh0lW16FNpPzO8vOFDWuF/HV/8uqrx8bHxPbYI9+xu4e5SgTHCZQ+EEadQHrpX/flazM2386JqDPJBuR6gDpytts4jjvETHQBXwHV8d9EV0ROE99CRAQRDByDMy+c/aiiTJvou01g0rlTA3KFA7qvRFFjXdGvWBWuNphpn+xciJjnAs9c/DIZizGIscgWg6BsnaD15SeJ9Y4FJF2THQ5azMgRSa/ty2oJwkmP/sjYjPJ8RCmk7hypmbAnMe4Co/IeId+pt5J2yO0pyNZoKnTRNB9pZoYNAV4Fwu5XjZDcOFoYjpNxOTm3/mH1u0+MLaEvHrQOu3+F8lx3OnJc9S9UubqZ7sycEcFyhkLg3LDKHwCPFdLfdDqLHDIFckBnzxw6f7SobJdL1mwgwRMzk0mjhgGW/vdyDDTrM7s7sF1J+VxG2eAcFtcMZ+cgPU9FL35hECS0djcgdFwVwDsK+LcU4Hy9A0Ym6dZJE4uk7J2L+01xk2cb397ngAaFP8Y1TaQhygSLAn+BMjRaaI60pvai333csJUYA3fSzpjJrZqKqiM3vyQNyqYzQmTrmu+Xd6P25gRUJED3U/BHJQVY7qeo+G12uXOzxLzaOTqmxmpq0g9GZvLaBCWNRBy213t1l7br9JPdb+kZPbi06hOebPaRHS+kHeAPRzXnA+PDcQbFOz5ei9e3SUwyfv/qjRB/F7JnPrcdTdKlkwvEDBEF+1hoTkxUg2bXTzgMwiV90Guwa7EC/6HzzS+Hd2J8euSx1cPOK5oRKFRJx0DUKadt4pRBhObrAVTCrpAocqOv/Qmb9RXZFtk+J3EH36DfKymc//H1qARUYareayyiHtj8ieM2bSHvobed8lorr5bGAFHjCm6woRwoWb8Fxonll2mm3IXirasKvgjqxVJ2Z8k2tvCWJZVV1n46mf2NkGI1JpHB1sTVvTSC0VLaeWh14NX4Zpnn8JoOggT9UZmO/GU0nLqFPo+mEz1x8n8vgSE71f3SXcI6DL9wmTCjFz2uQ0QDFtxOYNxEdukSxC+Bl5O7k2CcYrBu+KEnNFNxEHFK9MaHuZcQqw4ATOZ0tYO1iCwGQw9RqasLOL5rQvCDHwzXgSkNc0hyVijT+1HVHz0zte6d8y5x1b0OHKv4masaqcp+SqJ3H7+iYsm22O7R1rpVDqWQZ+Xdf7PPEAq2fTEknIZfmE6jqyKqlpFarEjHhfzfx3fj32IAXfB01yRFJBEccTi6hgFse3l5nWy9IqYLrxdRo31uop977npgkdGxH+ZcT2cOivL+ZXkbi8u3GN79J9lANJSqzvSoC6nC+A2SpXJOXjXC6jcUt/f/dtOArEwU6kE4sp+syM9oXZdCwHemGh5CZTAO+ItLSD7SSI/r4H4mSV9ZGN13dYJBZ1Zx4xN/s2Bkp45TJLxEnBFYTCcDSFHL76341o5NKLhftezjXEZx5pXQKEhPTDLrpTkaJmVYcz2etFtB68RIBrWvrpvKCj56JghjM3GIGaL/PhOc9kSBkOHt31eAQVn1dLieidd51fHpGeYdKuNAQ7BPn2rGzxpbfaasNbQxW5KNjGYeTz+SzwuSgpE/qiELw0m8tFTg8C9KPQFk6F8KKjL2f34lhratKor/9NfX/4kHs2kQbfZEpmaz9woUMHq6mSkNnspJ/o2fNh+v7LWKklxziNukvIKcEsD7efK581RfGx+swWNZUxddVnZMZ8dKi6fn08r8FAD0WweqVV896ZTEyXq7RS0Lqr4RHH6jgbpQj0szTzD5q109TrCoCyfofK2z0y5eVxchBPPh9YlBxtNUfffs/8bRnkzB4qze6AFhOe0xbAX4fiad5yEAohYXK9ScV+X4ifcgBVXJsY4nMhJ5sY5Pbmm+M/YFm0xjLFQZBH1HX/DElys9pEfoqPxgit2pMhsyRopqoLMOZAqP9z9txpTb9tSAXF1IUMXM8pB7MveaasX5J0Aosz8ac8yAS+HzBouwG/BHGIFIBIdLLUlq9n2mVWDgMAHoFOPD8Dj/wwOHbHq5tH6SX3FGecVoGCa3DmENPH689xcM9NEw1oshG5s8YpSHpANoFFo3MiGr5GwAT8TCNas0iA8diUrwVowZbLZ4fcNL5bAIYTdcJWophUbg5lo+bSGSn/FztJBo3LOILsq13I1HlS+rj6rUwAGtXcRMMcC/COO27a2OfOecu2WoccT40P1ujTtnKhhhiNebb/OJruVj0Gh1D89vjBLbBgNFsE0fl6G5WQrsGQePW5O+oamGMwX4x1FvnwGUFjwXASg/zUHmb+/CWnKmXjSbm+KW5YsBngJi/a91ssKmMJCaX8njTkKXOLvfntFl5IS84TMXnKX1j6TR6QA3iLaNTBOv4WRwPMby8PV4yYlZnUaSG86p09sU5rkTBcBA+IUlzREypbdiVtghW2rgPSPiYr3LiifVAxaG7DlXdlrfm/Q9cZOe6tUubSwkW8IbKVFzD5Wj4A7eq7O+qOHoM/ZPmKFzL9z/1sho/PL/k/fZq+G+jRNWUDKKGfuRVekzGQlz/peZUvFNWaH9GGvmvbgh5y9m7ogi0ba1SRWblFR8Z2+EEQAMp+kQISYH0vK1J3/VwuQ06VLOLMMduCYi5PxEhy5PDgVZBI060rUiI4E/qUfj9A="
	postData["__VIEWSTATEGENERATOR"] = "82F443D6"
	postData["__VIEWSTATEENCRYPTED"] = ""
	postData["__EVENTARGUMENT"] = ""
	postData["ddlQueryType"] = "byCity"
	postData["ddlCityList"] = "2"
	postData["ddlSchoolLevelByCity"] = "50"
	postData["ddlCourseTag"] = "安全教育"
	postData["Button1"] = "查詢"
	postData["__EVENTTARGET"] = ""
	payload := strings.NewReader(encodeSendData(postData))

	url := "https://www1.inservice.edu.tw/script/IndexQuery.aspx?city=2"
	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", "eeoq03njmzbwzz32lbifocvo,eeoq03njmzbwzz32lbifocvo; ASP.NET_SessionId=eeoq03njmzbwzz32lbifocvo")
	req.Header.Add("User-Agent", "PostmanRuntime/7.20.1")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Host", "www1.inservice.edu.tw")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Content-Length", "19788")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()
	// fmt.Println(string(body))

	// goquery 爬蟲取得資訊
	classes := []iClass{}
	dom, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	dom.Find(".cinfo-r1>tbody>tr>td:first-child").Each(func(i int, selection *goquery.Selection) {
		temp := iClass{ID: selection.Text(), Detail: postInserviceDetail(selection.Text())}
		classes = append(classes, temp)
	})

	dom.Find(".cinfo-r2").Each(func(i int, selection *goquery.Selection) {
		classes[i].Name = selection.Text()
	})

	dom.Find(".cinfo-r3").Each(func(i int, selection *goquery.Selection) {
		classes[i].Location = selection.Text()
	})

	nextPage := false
	dom.Find(".cssctsTitle2>td>a").Each(func(i int, selection *goquery.Selection) {
		href, ok := selection.Attr("href")
		if !ok {
			fmt.Println("error")
		}

		if regexp.MustCompile("ctl" + strconv.Itoa(page+1) + "'").Match([]byte(href)) {
			nextPage = true
		}
	})

	if nextPage {
		return append(classes, postInservicePage(page+1)...)
	}

	return classes
	// fmt.Println(res)
}

// encode map to string
func encodeSendData(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=%s&", key, url.QueryEscape(value))
	}
	return b.String()
}

// 取得詳細資料
func postInserviceDetail(id string) iDetail {
	url := "https://www1.inservice.edu.tw/NAPP/CourseView.aspx?cid=" + id

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("User-Agent", "PostmanRuntime/7.16.3")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Host", "www1.inservice.edu.tw")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()

	dom, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// goquery 爬蟲取得資訊
	var detail iDetail
	// 報名狀態
	dom.Find("#ctl00_CPH_Content_pl_courseData > div:nth-child(7) > table > tbody > tr > td:nth-child(2) > table > tbody > tr:nth-child(1) > td.cw_table_info.c4.cs1").Each(func(i int, selection *goquery.Selection) {
		detail.SignUpStatus = selection.Text()
	})

	// 報名時間
	dom.Find("#ctl00_CPH_Content_pl_courseData > div:nth-child(7) > table > tbody > tr > td:nth-child(2) > table > tbody > tr:nth-child(1) > td.cw_table_info.c2.cs1").Each(func(i int, selection *goquery.Selection) {
		detail.SignUpTime = selection.Text()
	})

	// 上課日期
	dom.Find("#ctl00_CPH_Content_pl_courseData > div:nth-child(4) > table > tbody > tr > td:nth-child(2) > table > tbody > tr:nth-child(1) > td.cw_table_info.c2.cs1").Each(func(i int, selection *goquery.Selection) {
		detail.AttendClassTime = selection.Text()
	})

	// 研習時數
	dom.Find("#ctl00_CPH_Content_pl_courseData > div:nth-child(4) > table > tbody > tr > td:nth-child(2) > table > tbody > tr:nth-child(3) > td.cw_table_info.c2.cs3").Each(func(i int, selection *goquery.Selection) {
		detail.StudyHours = selection.Text()
	})

	// 開課地點
	dom.Find("#ctl00_CPH_Content_pl_courseData > div:nth-child(4) > table > tbody > tr > td:nth-child(2) > table > tbody > tr:nth-child(2) > td.cw_table_info.c2.cs3").Each(func(i int, selection *goquery.Selection) {
		detail.Location = selection.Text()
		detail.MapDetail = getMapDuration(selection.Text())
	})

	return detail
}

// 取得map資料
func getMapDuration(destinations string) gMap {
	// 從config取得map key
	mapConfig, err := config.NewConfig("ini", "conf/env.conf")
	mapKey := mapConfig.String("gMapKey")

	sendData := make(map[string]string)
	sendData["units"] = "imperial"
	sendData["origins"] = "242新北市新莊區中正路893巷120號"
	sendData["destinations"] = destinations
	sendData["mode"] = "transit"
	sendData["key"] = mapKey
	url := "https://maps.googleapis.com/maps/api/distancematrix/json?" + encodeSendData(sendData)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	// json string to struct
	g := gMap{}
	err2 := json.Unmarshal(body, &g)
	if err2 != nil {
		log.Fatalln(err2)
	}

	return g
}
