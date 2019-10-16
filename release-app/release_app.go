package main

import (
	"context"
	"flag"
	"github.com/go-openapi/strfmt"
	"io/ioutil"
	"openpitrix.io/openpitrix/pkg/client/app"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"strings"
)

func main() {
	var path string
	flag.StringVar(&path, "path", "./package/", "need package path.eg./your/path/to/pkg/")
	flag.Parse()

	fileInfoList, err := ioutil.ReadDir(path)
	if err != nil {
		logger.Error(nil, "read dir path error: %s",err.Error())
	}

	for _, f := range fileInfoList {
		filePath := path + f.Name()
		var appName string
		segName := strings.Split(f.Name(), "-")
		if len(segName) > 0 {
			appName = segName[0]
		}
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			logger.Error(nil, "")
		}
		pkg := strfmt.Base64(content)

		client, err := app.NewAppManagerClient()
		if err != nil {
			panic(err)
		}

		createReq := &pb.CreateAppRequest{
			VersionPackage: pbutil.ToProtoBytes(pkg),
			Name:    pbutil.ToProtoString(appName),
			VersionType:    pbutil.ToProtoString("helm"),
		}

		res, err := client.CreateApp(context.Background(), createReq)
		if err != nil {
			logger.Error(nil, "create app error: %s", err.Error())
		} else {
			submitReq := &pb.SubmitAppVersionRequest{
				VersionId: res.VersionId,
			}
			_, err = client.SubmitAppVersion(context.Background(), submitReq)
			if err != nil {
				logger.Error(nil, "submit app error: %s", err.Error())
			}

			passReq := &pb.PassAppVersionRequest{
				VersionId: res.VersionId,
			}
			_, err = client.AdminPassAppVersion(context.Background(), passReq)
			if err != nil {
				logger.Error(nil, "pass app error: %s", err.Error())
			}

			releaseReq := &pb.ReleaseAppVersionRequest{
				VersionId: res.VersionId,
			}
			_, err = client.ReleaseAppVersion(context.Background(), releaseReq)
			if err != nil {
				logger.Error(nil, "release app error: %s", err.Error())
			}
		}
	}

}
