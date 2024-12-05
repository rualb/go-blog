package e2e

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	xcmd "go-blog/internal/cmd"
	"go-blog/internal/config"
	"go-blog/internal/config/consts"
	"go-blog/internal/repository"
	"go-blog/internal/service"
	"go-blog/internal/util/utilhttp"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func setup() {

	cfgSrc := config.MustNewAppConfigSource()
	cfg := cfgSrc.Config()
	db := repository.MustNewRepository(cfg) // not via app-service

	{
		// migrate
		for _, x := range []any{
			&service.UserAccount{},
			&service.VaultKey{},
		} {
			if err := db.AutoMigrate(x); err != nil {
				log.Fatalf("Migration: %v", err)
			}
		}

		// seed
		{
			key := &service.VaultKey{}
			key.ID = "test-only-not-a-secret-key"
			key.AuthKey = base64.StdEncoding.EncodeToString([]byte(strings.Repeat("A", 64)))
			// insert or update
			if res := db.Save(key); res.Error != nil {
				log.Fatalf("New vault key: %v", res.Error)
			}
		}
	}
}

func TestMain(m *testing.M) {
	fmt.Println("TestMain")
	setup()
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestBlogAddData(t *testing.T) {

	os.Setenv("APP_ENV", "testing")
	os.Setenv("APP_BLOG_ADMIN", "true")

	srv := service.MustNewAppServiceProd()
	blogSrv := srv.Blog()
	dao := blogSrv.Posts()

	for i := 1; i < 199; i++ {
		tmp := service.BlogPost{
			Code:            fmt.Sprintf("test-post-%v", i),
			Title:           fmt.Sprintf("Test post %v", i),
			ContentMarkdown: fmt.Sprintf("Test post %v **Content**", i),
			// ContentHTML: fmt.Sprintf("Test post %v <b>Content</b>", i),
		}

		if id, _ := dao.Code(tmp.Code); id > 0 {
			continue
		}

		tmp.Fill()
		if err := dao.Create(&tmp); err != nil {
			t.Fatal(err)
		}
	}

}

// TestHealthController_Check_Stats tests the ?cmd=stats case in the Check method
func TestBlog(t *testing.T) {
	// Setup Echo context

	//   Trafalgar

	os.Setenv("APP_ENV", "testing")

	cmd := xcmd.Command{}

	go cmd.Exec()

	time.Sleep(1 * time.Second)

	{

		data := &service.BlogPost{
			Code:            "test-post-1",
			Title:           "Test post 1",
			ContentMarkdown: "Test blog 1 **Content**",
			// ContentHTML: "Test post 1 <b>Content</b>",
		}

		data.Fill()

		srv := cmd.AppService
		blogSrv := srv.Blog()
		repoSrv := srv.Repository()

		if resp := repoSrv.Where("code = ?", data.Code).Delete(&service.BlogPost{}); resp.Error != nil {
			t.Fatal(resp.Error)
		}

		if dataTmp, err := blogSrv.Posts().FindByCode("test-post-1"); err != nil {
			t.Fatal(err)
		} else if dataTmp == nil {

			res := repoSrv.Create(data)

			if res.Error != nil {
				t.Fatal(res.Error)
			}

		}

		urls := []struct {
			title  string
			url    string
			query  map[string]string
			search []string
		}{

			{title: "test 1", search: []string{
				data.ContentMarkdown,
				// data.ContentHTML,
			},
				url:   "http://127.0.0.1:38180" + strings.ReplaceAll(consts.PathBlogPostsEntityByCodeAPI, ":code", data.Code),
				query: map[string]string{}},
		}

		for _, itm := range urls {

			t.Run(itm.title, func(t *testing.T) {

				t.Logf("url %v", itm.url)
				dataArr, err := utilhttp.GetBytes(itm.url, itm.query, nil)

				if err != nil {
					t.Errorf("Error : %v", err)
				}
				dataTxt := string(dataArr)
				for _, v := range itm.search {

					d, _ := json.Marshal(&v)
					if !strings.Contains(string(dataTxt), string(d)) {
						t.Errorf("Error on %v", itm.url)
					}
				}

			})

		}

	}

	cmd.Stop()

	time.Sleep(1 * time.Second)

}
