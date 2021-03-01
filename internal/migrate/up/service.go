package up

import (
	"encoding/json"
	"esctl/internal/migrate"
	"esctl/pkg/es"
	"esctl/pkg/log"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type IService interface {
	InitMigrationHistoryRepo() error
	SaveMigrationHistoryEntry(name string) error
	DeleteMigrationHistoryEntry(name string) error
	GetLastMigrationName() (string, error)
}

type service struct {
	logHelper log.IHelper
	esHelper  es.IHelper
}

func NewService(logHelper log.IHelper, esHelper es.IHelper) IService {
	return &service{
		logHelper,
		esHelper,
	}
}

// 初始化迁移历史仓库
func (s *service) InitMigrationHistoryRepo() error {
	esIndexName := migrate.MIGRATION_HISTORY_ES_INDEX_NAME

	listResp, err := s.esHelper.ListIndices(esIndexName)
	if err != nil {
		return err
	}

	if len(listResp.Items) > 0 {
		return nil
	}

	esIndexBody := migrate.MIGRATION_HISTORY_ES_INDEX_BODY
	_, err = s.esHelper.CreateIndex(esIndexName, esIndexBody)
	if err != nil {
		return err
	}

	return nil
}

// 保存迁移历史条目
func (s *service) SaveMigrationHistoryEntry(name string) error {
	esIndexName := migrate.MIGRATION_HISTORY_ES_INDEX_NAME
	esDocID := name
	esDocBody := []byte(fmt.Sprintf(`{"name":"%s"}`, name))
	if err := s.esHelper.SaveDoc(esIndexName, esDocID, esDocBody); err != nil {
		return err
	}

	return nil
}

// 删除迁移历史条目
func (s *service) DeleteMigrationHistoryEntry(name string) error {
	esIndexName := migrate.MIGRATION_HISTORY_ES_INDEX_NAME
	esDocID := name
	if err := s.esHelper.DeleteDoc(esIndexName, esDocID); err != nil {
		return err
	}

	return nil
}

// 获取最后迁移名称
func (s *service) GetLastMigrationName() (string, error) {
	esIndexName := migrate.MIGRATION_HISTORY_ES_INDEX_NAME
	esSearchBody := []byte(`{"sort":[{"_id":{"order":"desc"}}],"size":1}`)
	listResp, err := s.esHelper.ListDocs(esIndexName, esSearchBody)
	if err != nil {
		return "", err
	}

	if len(listResp.Hits.Hits) == 0 {
		return "", nil
	}

	result := &migrate.MigrationHistoryEntry{}
	if err := json.Unmarshal(listResp.Hits.Hits[0].Source, result); err != nil {
		return "", err
	}

	return result.Name, nil
}

// 列出向上迁移文件名（升序）
func (s *service) listMigrateUpFileNames(dir string) ([]string, error) {
	fd, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	files, err := fd.Readdir(-1)
	if err != nil {
		return nil, err
	}

	var res []string
	for _, file := range files {
		fName := file.Name()
		if strings.HasSuffix(fName, migrate.MIGRATION_UP_FILE_SUFFIX) {
			res = append(res, fName)
		}
	}

	sort.Strings(res)

	return res, nil
}

func (s *service) parseMigrationFile(file string) (*migrate.Migration, error) {
	viper.SetConfigFile(file)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "failed to read migration file")
	}

	res := &migrate.Migration{}
	if err := viper.Unmarshal(res); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal migration file content")
	}

	return res, nil
}
