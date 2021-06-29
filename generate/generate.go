package generate

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"
	"time"

	"CharlesBases/mysql-gen-go/utils"
)

type Infor struct {
	Config *utils.GlobalConfig
	Struct *utils.Struct
}

func (infor *Infor) GenModel(wr io.Writer) {
	temp := template.New(infor.Struct.StructName)
	temp.Funcs(template.FuncMap{
		"package": func() string {
			return infor.Config.Package
		},
		"imports": func() template.HTML {
			importsbuilder := strings.Builder{}
			for _, field := range *infor.Struct.Fields {
				if strings.HasSuffix(field.Type, "time.Time") {
					importsbuilder.WriteString(fmt.Sprintf("%s\n\t", `"time"`))
				}
			}
			return template.HTML(importsbuilder.String())
		},
		"html": func(source string) template.HTML {
			return template.HTML(source)
		},
		"gormDB": func() template.HTML {
			strBuilder := strings.Builder{}
			strBuilder.WriteString("db *gorm.DB " + "`" + `json:"-" gorm:"-"` + "`")
			return template.HTML(strBuilder.String())
		},
	})
	modelTemplate, err := temp.Parse(modeltemplate)
	if err != nil {
		fmt.Print(fmt.Sprintf("[%s]----------%c[%d;%d;%dmgen model error: %s%c[0m\n", time.Now().Format("2006-01-02 15:04:05"), 0x1B, 0 /*字体*/, 0 /*背景*/, 31 /*前景*/, err.Error(), 0x1B))
		os.Exit(1)
	}
	modelTemplate.Execute(wr, infor.Struct)
}