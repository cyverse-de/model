package model

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/cyverse-de/configurate"

	"github.com/spf13/viper"
)

func JSONData(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	c, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return c, err
}

var (
	s   *Job
	cfg *viper.Viper
)

func _initconfig(t *testing.T) {
	var err error
	cfg, err = configurate.InitDefaults("test/test_config.yaml", configurate.JobServicesDefaults)
	if err != nil {
		t.Fatal(err)
	}
	cfg.Set("condor.run_on_nfs", true)
	cfg.Set("irods.base", "/path/to/irodsbase")
	cfg.Set("irods.host", "hostname")
	cfg.Set("irods.port", "1247")
	cfg.Set("irods.user", "user")
	cfg.Set("irods.pass", "pass")
	cfg.Set("irods.zone", "test")
	cfg.Set("irods.resc", "")
	cfg.Set("condor.log_path", "/path/to/logs")
	cfg.Set("condor.porklock_tag", "test")
	cfg.Set("condor.filter_files", "foo,bar,baz,blippy")
}

func inittestsFile(t *testing.T, filename string) *Job {
	var err error
	if cfg == nil {
		_initconfig(t)
	}
	data, err := JSONData(filename)
	if err != nil {
		t.Fatal(err)
	}
	submission, err := NewFromData(cfg, data)
	if err != nil {
		t.Fatal(err)
	}
	return submission
}

func _inittests(t *testing.T, memoize bool) *Job {
	if s == nil || !memoize {
		s = inittestsFile(t, "test/test_submission.json")
	}
	return s
}

func inittests(t *testing.T) *Job {
	return _inittests(t, true)
}

func TestJSONParsing(t *testing.T) {
	inittests(t)
}

func TestNaivelyQuote(t *testing.T) {
	test1 := naivelyquote("foo")
	test2 := naivelyquote("'foo'")
	test3 := naivelyquote("foo'oo")
	test4 := naivelyquote("'foo'oo'")
	test5 := naivelyquote("foo''oo")
	test6 := naivelyquote("'foo''oo'")
	test7 := naivelyquote("f'oo'oo")
	test8 := naivelyquote("'f'oo'oo'")

	if test1 != "'foo'" {
		t.Errorf("naivelyquote returned %s instead of 'foo'", test1)
	}
	if test2 != "'''foo'''" {
		t.Errorf("naivelyquote returned %s instead of '''foo'''", test2)
	}
	if test3 != "'foo''oo'" {
		t.Errorf("naivelyquote returned %s instead of 'foo''oo'", test3)
	}
	if test4 != "'''foo''oo'''" {
		t.Errorf("naivelyquote returned %s instead of '''foo''oo'''", test4)
	}
	if test5 != "'foo''''oo'" {
		t.Errorf("naivelyquote returned %s instead of 'foo''''oo'", test5)
	}
	if test6 != "'''foo''''oo'''" {
		t.Errorf("naivelyquote returned %s instead of '''foo''''oo'''", test6)
	}
	if test7 != "'f''oo''oo'" {
		t.Errorf("naivelyquote returned %s instead of 'f''oo''oo'", test7)
	}
	if test8 != "'''f''oo''oo'''" {
		t.Errorf("naivelyquote returned %s instead of '''f''oo''oo'''", test8)
	}
}

func TestQuote(t *testing.T) {
	test1 := quote("foo")
	test2 := quote("'foo'")
	test3 := quote("foo'oo")
	test4 := quote("'foo'oo'")
	test5 := quote("foo''oo")
	test6 := quote("'foo''oo'")
	test7 := quote("f'oo'oo")
	test8 := quote("'f'oo'oo'")

	if test1 != "'foo'" {
		t.Errorf("naivelyquote returned %s instead of 'foo'", test1)
	}
	if test2 != "'foo'" {
		t.Errorf("naivelyquote returned %s instead of 'foo'", test2)
	}
	if test3 != "'foo''oo'" {
		t.Errorf("naivelyquote returned %s instead of 'foo''oo'", test3)
	}
	if test4 != "'foo''oo'" {
		t.Errorf("naivelyquote returned %s instead of 'foo''oo'", test4)
	}
	if test5 != "'foo''''oo'" {
		t.Errorf("naivelyquote returned %s instead of 'foo''''oo'", test5)
	}
	if test6 != "'foo''''oo'" {
		t.Errorf("naivelyquote returned %s instead of 'foo''''oo'", test6)
	}
	if test7 != "'f''oo''oo'" {
		t.Errorf("naivelyquote returned %s instead of 'f''oo''oo'", test7)
	}
	if test8 != "'f''oo''oo'" {
		t.Errorf("naivelyquote returned %s instead of 'f''oo''oo'", test8)
	}
}

func TestIRODSBase(t *testing.T) {
	s := inittests(t)
	if s.IRODSBase != "/path/to/irodsbase" {
		t.Errorf("The IRODS base directory was set to '%s' instead of '/path/to/irodsbase'", s.IRODSBase)
	}
}

func TestDescription(t *testing.T) {
	s := inittests(t)
	if s.Description != "this is a description" {
		t.Errorf("The description was '%s' instead of 'this is a description'", s.Description)
	}
}

func TestEmail(t *testing.T) {
	s := inittests(t)
	if s.Email != "wregglej@iplantcollaborative.org" {
		t.Errorf("The email was '%s' instead of wregglej@iplantcollaborative.org", s.Email)
	}
}

func TestName(t *testing.T) {
	s := inittests(t)
	if s.Name != "Word_Count_analysis1__" {
		t.Errorf("The name field was '%s' instead of 'Word_Count_analysis1__'", s.Name)
	}
}

func TestUsername(t *testing.T) {
	s := inittests(t)
	if s.Submitter != "test_this_is_a_test" {
		t.Errorf("The username was '%s' instead of 'test_this_is_a_test'", s.Submitter)
	}
}

func TestAppID(t *testing.T) {
	s := inittests(t)
	if s.AppID != "c7f05682-23c8-4182-b9a2-e09650a5f49b" {
		t.Errorf("The app_id was '%s' instead of 'c7f05682-23c8-4182-b9a2-e09650a5f49b'", s.AppID)
	}
}

func TestCreateOutputSubdir(t *testing.T) {
	s := inittests(t)
	if !s.CreateOutputSubdir {
		t.Errorf("create_output_subdir was false when it should have been true")
	}
}

func TestRequestType(t *testing.T) {
	s := inittests(t)
	if s.RequestType != "submit" {
		t.Errorf("request_type was '%s' instead of 'submit'", s.RequestType)
	}
}

func TestAppDescription(t *testing.T) {
	s := inittests(t)
	if s.AppDescription != "this is an app description" {
		t.Errorf("app_description was '%s' instead of 'this is an app description'", s.AppDescription)
	}
}

func TestOutputDir(t *testing.T) {
	s := inittests(t)
	if s.OutputDir != "/iplant/home/wregglej/analyses/Word_Count_analysis1-2015-09-17-21-42-20.9" {
		t.Errorf("output_dir was '%s' instead of '/iplant/home/wregglej/analyses/Word_Count_analysis1-2015-09-17-21-42-20.9'", s.OutputDir)
	}
}

func TestWikiURL(t *testing.T) {
	s := inittests(t)
	if s.WikiURL != "https://pods.iplantcollaborative.org/wiki/display/DEapps/WordCount" {
		t.Errorf("wiki_url was '%s' instead of 'https://pods.iplantcollaborative.org/wiki/display/DEapps/WordCount'", s.WikiURL)
	}
}

func TestUUID(t *testing.T) {
	s := inittests(t)
	if s.InvocationID != "07b04ce2-7757-4b21-9e15-0b4c2f44be26" {
		t.Errorf("uuid was '%s' instead of '07b04ce2-7757-4b21-9e15-0b4c2f44be26'", s.InvocationID)
	}
}

func TestNotify(t *testing.T) {
	s := inittests(t)
	if !s.Notify {
		t.Errorf("notify was false instead of true")
	}
}

func TestExecutionTarget(t *testing.T) {
	s := inittests(t)
	if s.ExecutionTarget != "condor" {
		t.Errorf("execution_target was '%s' instead of 'condor'", s.ExecutionTarget)
	}
}

func TestAppName(t *testing.T) {
	s := inittests(t)
	if s.AppName != "Word Count" {
		t.Errorf("app_name was '%s' instead of 'Word Count'", s.AppName)
	}
}

func TestUserGroups(t *testing.T) {
	s := _inittests(t, false)
	expected := `{"groups:foo","groups:bar","groups:baz"}`
	actual := s.FormatUserGroups()
	if actual != expected {
		t.Errorf("FormatUserGroups() returned '%s' when it should have returned '%s'", actual, expected)
	}
}

func TestEmptyUserGroups(t *testing.T) {
	s := _inittests(t, false)
	s.UserGroups = []string{}
	expected := `{}`
	actual := s.FormatUserGroups()
	if actual != expected {
		t.Errorf("FormatUserGroups() returned '%s' when it should have returned '%s'", actual, expected)
	}
}

func TestNilUserGroups(t *testing.T) {
	s := _inittests(t, false)
	s.UserGroups = []string(nil)
	expected := `{}`
	actual := s.FormatUserGroups()
	if actual != expected {
		t.Errorf("FormatUserGroups() returned '%s' when it should have returned '%s'", actual, expected)
	}
}

func TestMissingUserGroups(t *testing.T) {
	s := inittestsFile(t, "test/no_groups_submission.json")
	expected := `{}`
	actual := s.FormatUserGroups()
	if actual != expected {
		t.Errorf("FormatUserGroups() returned '%s' when it should have returned '%s'", actual, expected)
	}
}

func TestDirname(t *testing.T) {
	s := _inittests(t, false)
	s.NowDate = time.Now().Format(nowfmt)
	expected := fmt.Sprintf("%s-%s", s.Name, s.NowDate)
	actual := s.DirectoryName()
	if actual != expected {
		t.Errorf("Dirname() returned '%s' when it should have returned '%s'", actual, expected)
	}
}

func TestCondorLogDir(t *testing.T) {
	s := _inittests(t, false)
	s.NowDate = time.Now().Format(nowfmt)
	logPath := cfg.GetString("condor.log_path")
	expected := fmt.Sprintf("%s/", path.Join(logPath, s.Submitter, s.DirectoryName()))
	actual := s.CondorLogDirectory()
	if actual != expected {
		t.Errorf("CondorLogDir() returned '%s' when it should have returned '%s'", actual, expected)
	}
}

func TestIRODSConfig(t *testing.T) {
	s := _inittests(t, false)
	s.NowDate = time.Now().Format(nowfmt)
	expected := path.Join("logs", "irods-config")
	actual := s.IRODSConfig()
	if actual != expected {
		t.Errorf("IRODSConfig() returned '%s' when it should have returned '%s'", actual, expected)
	}
}

func TestOutputDirectory1(t *testing.T) {
	s := _inittests(t, false)
	s.OutputDir = ""
	expected := path.Join(s.IRODSBase, s.Submitter, "analyses", s.DirectoryName())
	actual := s.OutputDirectory()
	if actual != expected {
		t.Errorf("OutputDirectory() returned '%s' when it should have returned '%s'", actual, expected)
	}
}

func TestOutputDirectory2(t *testing.T) {
	s := _inittests(t, false)
	expected := path.Join(s.OutputDir, s.DirectoryName())
	actual := s.OutputDirectory()
	if actual != expected {
		t.Errorf("OutputDirectory() returned '%s' when it should have returned '%s'", actual, expected)
	}
}

func TestOutputDirectory3(t *testing.T) {
	s := _inittests(t, false)
	s.CreateOutputSubdir = false
	expected := strings.TrimSuffix(s.OutputDir, "/")
	actual := s.OutputDirectory()
	if actual != expected {
		t.Errorf("OutputDirectory() returned '%s' when it should have returned '%s'", actual, expected)
	}
}

func TestDataContainers(t *testing.T) {
	s := inittests(t)
	dc := s.DataContainers()
	dclen := len(dc)
	if dclen != 2 {
		t.Errorf("The number of data containers was '%d' instead of 2", dclen)
	}

	vfs := dc[0]
	if vfs.Name != "vf-name1" {
		t.Errorf("The VolumesFrom name was '%s' when it should have been 'vf-name1'", vfs.Name)
	}
	if vfs.NamePrefix != "vf-prefix1" {
		t.Errorf("The VolumesFrom prefix was '%s' when it should have been 'vf-prefix1'", vfs.NamePrefix)
	}
	if vfs.Tag != "vf-tag1" {
		t.Errorf("The VolumesFrom tag was '%s' when it should have been 'vf-tag1'", vfs.Tag)
	}
	if vfs.URL != "vf-url1" {
		t.Errorf("The VolumesFrom url was '%s' when it should have been 'vf-url1'", vfs.URL)
	}
	if vfs.HostPath != "/host/path1" {
		t.Errorf("The VolumesFrom host path was '%s' when it should have been '/host/path1'", vfs.HostPath)
	}
	if vfs.ContainerPath != "/container/path1" {
		t.Errorf("The VolumesFrom container path was '%s' when it should have been '/container/path1'", vfs.ContainerPath)
	}
	if !vfs.ReadOnly {
		t.Error("The VolumesFrom read-only field was false when it should have been true.")
	}

	vfs = dc[1]
	if vfs.Name != "vf-name2" {
		t.Errorf("The VolumesFrom name was '%s' when it should have been 'vf-name2'", vfs.Name)
	}
	if vfs.NamePrefix != "vf-prefix2" {
		t.Errorf("The VolumesFrom prefix was '%s' when it should have been 'vf-prefix2'", vfs.NamePrefix)
	}
	if vfs.Tag != "vf-tag2" {
		t.Errorf("The VolumesFrom tag was '%s' when it should have been 'vf-tag2'", vfs.Tag)
	}
	if vfs.URL != "vf-url2" {
		t.Errorf("The VolumesFrom url was '%s' when it should have been 'vf-url2'", vfs.URL)
	}
	if vfs.HostPath != "/host/path2" {
		t.Errorf("The VolumesFrom host path was '%s' when it should have been '/host/path2'", vfs.HostPath)
	}
	if vfs.ContainerPath != "/container/path2" {
		t.Errorf("The VolumesFrom container path was '%s' when it should have been '/container/path2'", vfs.ContainerPath)
	}
	if !vfs.ReadOnly {
		t.Error("The VolumesFrom read-only field was false when it should have been true.")
	}
}

func validateContainerImage(t *testing.T, actualImage, expectedImage ContainerImage) {
	actual := actualImage.ID
	expected := expectedImage.ID
	if actual != expected {
		t.Errorf("ID was %s instead of %s", actual, expected)
	}
	actual = actualImage.Name
	expected = expectedImage.Name
	if actual != expected {
		t.Errorf("Name was %s instead of %s", actual, expected)
	}
	actual = actualImage.Tag
	expected = expectedImage.Tag
	if actual != expected {
		t.Errorf("Tag was %s instead of %s", actual, expected)
	}
	actual = actualImage.URL
	expected = expectedImage.URL
	if actual != expected {
		t.Errorf("URL was %s instead of %s", actual, expected)
	}
	actual = actualImage.OSGImagePath
	expected = expectedImage.OSGImagePath
	if actual != expected {
		t.Errorf("OSGImagePath was %s instead of %s", actual, expected)
	}
}

func min(n int, ns ...int) int {
	minValue := n
	for _, num := range ns {
		if num < minValue {
			minValue = num
		}
	}
	return minValue
}

func validateContainerImages(t *testing.T, actualImages, expectedImages []ContainerImage) {
	actuallen := len(actualImages)
	expectedlen := len(expectedImages)
	if actuallen != expectedlen {
		t.Errorf("ContainerImages() returned %d ContainerImages instead of %d", actuallen, expectedlen)
	}
	for i := 0; i < min(actuallen, expectedlen); i++ {
		validateContainerImage(t, actualImages[i], expectedImages[i])
	}
}

func TestContainerImages(t *testing.T) {
	s := inittests(t)
	actualImages := s.ContainerImages()
	expectedImages := []ContainerImage{
		{
			ID:           "fc210a84-f7cd-4067-939c-a68ec3e3bd2b",
			Name:         "gims.iplantcollaborative.org:5000/backwards-compat",
			Tag:          "latest",
			URL:          "https://registry.hub.docker.com/u/discoenv/backwards-compat",
			OSGImagePath: "",
		},
	}
	validateContainerImages(t, actualImages, expectedImages)
}

func TestContainerImagesOSG(t *testing.T) {
	s := inittestsFile(t, "test/test_submission_osg.json")
	actualImages := s.ContainerImages()
	expectedImages := []ContainerImage{
		{
			ID:           "fc210a84-f7cd-4067-939c-a68ec3e3bd2b",
			Name:         "gims.iplantcollaborative.org:5000/backwards-compat",
			Tag:          "latest",
			URL:          "https://registry.hub.docker.com/u/discoenv/backwards-compat",
			OSGImagePath: "/path/to/image",
		},
	}
	validateContainerImages(t, actualImages, expectedImages)
}

func TestFileMetadata(t *testing.T) {
	s := inittests(t)
	fm := s.FileMetadata
	actual := fm[0].Attribute
	expected := "attr1"
	if actual != expected {
		t.Errorf("Attribute was %s instead of %s", actual, expected)
	}
	actual = fm[0].Value
	expected = "value1"
	if actual != expected {
		t.Errorf("Value was %s instead of %s", actual, expected)
	}
	actual = fm[0].Unit
	expected = "unit1"
	if actual != expected {
		t.Errorf("Unit was %s instead of %s", actual, expected)
	}
	actual = fm[1].Attribute
	expected = "attr2"
	if actual != expected {
		t.Errorf("Attribute was %s instead of %s", actual, expected)
	}
	actual = fm[1].Value
	expected = "value2"
	if actual != expected {
		t.Errorf("Value was %s instead of %s", actual, expected)
	}
	actual = fm[1].Unit
	expected = "unit2"
	if actual != expected {
		t.Errorf("Unit was %s instead of %s", actual, expected)
	}
}

func TestFileMetadataArgument(t *testing.T) {
	s := inittests(t)
	fm := s.FileMetadata
	actual := fm[0].Argument()
	expected := []string{"-m", "attr1,value1,unit1"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Argument() returned %#v instead of %#v", actual, expected)
	}
	actual = fm[1].Argument()
	expected = []string{"-m", "attr2,value2,unit2"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Argument() returned %#v instead of %#v", actual, expected)
	}
}

func TestSubmissionFileMetadataArguments(t *testing.T) {
	s := inittests(t)
	actual := MetadataArgs(s.FileMetadata).FileMetadataArguments()
	expected := []string{
		"-m", "attr1,value1,unit1",
		"-m", "attr2,value2,unit2",
		"-m", "ipc-analysis-id,c7f05682-23c8-4182-b9a2-e09650a5f49b,UUID",
		"-m", "ipc-execution-id,07b04ce2-7757-4b21-9e15-0b4c2f44be26,UUID",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("FileMetadataArguments() returned %#v instead of %#v", actual, expected)
	}
}

func TestInputs(t *testing.T) {
	s := inittests(t)
	inputs := s.Inputs()
	actual := len(inputs)
	expected := 2
	if actual != expected {
		t.Errorf("Number of inputs was %d instead of %d", actual, expected)
	}
}

func TestOutputs(t *testing.T) {
	s := inittests(t)
	outputs := s.Outputs()
	actual := len(outputs)
	expected := 2
	if actual != expected {
		t.Errorf("Number of outputs was %d instead of %d", actual, expected)
	}
}

func TestExcludeArguments(t *testing.T) {
	s := inittests(t)
	actual := s.ExcludeArguments()
	expected := []string{"foo", "bar", "baz", "blippy"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("ExcludeArguments() returned:\n\t%#vinstead of:\n\t%#v", actual, expected)
	}
	s.Steps[0].Config.Inputs[0].Retain = false
	actual = s.ExcludeArguments()
	expected = []string{"Acer-tree.txt", "foo", "bar", "baz", "blippy"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("ExcludeArguments() returned:\n\t%sinstead of:\n\t%s", actual, expected)
	}
	s.Steps[0].Config.Outputs[1].Retain = false
	actual = s.ExcludeArguments()
	expected = []string{"Acer-tree.txt", "/de-app-work/logs/", "foo", "bar", "baz", "blippy"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("ExcludeArguments() returned:\n\t%sinstead of:\n\t%s", actual, expected)
	}
	s.ArchiveLogs = false
	actual = s.ExcludeArguments()
	expected = []string{"Acer-tree.txt", "/de-app-work/logs/", "foo", "bar", "baz", "blippy", "logs"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("ExcludeArguments() returned:\n\t%sinstead of:\n\t%s", actual, expected)
	}
	_inittests(t, false)
}

func TestAddRequiredMetadata(t *testing.T) {
	s := inittests(t)
	found := false
	var a FileMetadata
	for _, md := range s.FileMetadata {
		if md.Attribute == "ipc-analysis-id" {
			found = true
			a = md
		}
	}
	if !found {
		t.Errorf("ipc-analysis-id was not in the file metadata before AddRequiredMetadata() was called")
	}
	if found {
		if a.Value != s.AppID {
			t.Errorf("Value was set to %s instead of %s", a.Value, s.AppID)
		}
		if a.Unit != "UUID" {
			t.Errorf("Unit was set to %s instead of %s", a.Unit, "UUID")
		}
	}
	found = false
	var e FileMetadata
	for _, md := range s.FileMetadata {
		if md.Attribute == "ipc-execution-id" {
			found = true
			e = md
		}
	}
	if !found {
		t.Errorf("ipc-execution-id was not in the file metadata before AddRequiredMetadata() was called")
	}
	if found {
		if e.Value != s.InvocationID {
			t.Errorf("Value was set to %s instead of %s", e.Value, s.InvocationID)
		}
		if e.Unit != "UUID" {
			t.Errorf("Unit was set to %s instead of %s", e.Unit, "UUID")
		}
	}
	_inittests(t, false)
}

func TestFinalOutputArguments(t *testing.T) {
	s := inittests(t)
	s.AddRequiredMetadata()
	actual := s.FinalOutputArguments("exclude.txt")
	outputdir := s.OutputDirectory()
	expected := []string{
		"put",
		"--user", "test_this_is_a_test",
		"--destination", fmt.Sprintf("%s", outputdir),
		"--config", "/configs/irods-config",
		"-m", "attr1,value1,unit1",
		"-m", "attr2,value2,unit2",
		"-m", "ipc-analysis-id,c7f05682-23c8-4182-b9a2-e09650a5f49b,UUID",
		"-m", "ipc-execution-id,07b04ce2-7757-4b21-9e15-0b4c2f44be26,UUID",
		"--exclude", "exclude.txt",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("FinalOutputArguments() returned:\n\t%#v\ninstead of:\n\t%#v", actual, expected)
	}
	s.SkipParentMetadata = true
	actual = s.FinalOutputArguments("exclude.txt")
	expected = []string{
		"put",
		"--user", "test_this_is_a_test",
		"--destination", fmt.Sprintf("%s", outputdir),
		"--config", "/configs/irods-config",
		"-m", "attr1,value1,unit1",
		"-m", "attr2,value2,unit2",
		"-m", "ipc-analysis-id,c7f05682-23c8-4182-b9a2-e09650a5f49b,UUID",
		"-m", "ipc-execution-id,07b04ce2-7757-4b21-9e15-0b4c2f44be26,UUID",
		"--exclude", "exclude.txt",
		"--skip-parent-meta",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("FinalOutputArguments() returned:\n\t%#v\ninstead of:\n\t%#v", actual, expected)
	}
	_inittests(t, false)
}

func TestCPURequest(t *testing.T) {
	s := inittests(t)
	cpu := s.CPURequest()
	var expected float32
	expected = 0
	if cpu != expected {
		t.Errorf("CPU request was %f, not %f", cpu, expected)
	}
}

func TestMemoryRequest(t *testing.T) {
	s := inittests(t)
	mem := s.MemoryRequest()
	var expected int64
	expected = 2048
	if mem != expected {
		t.Errorf("Memory request was %d, not %d", mem, expected)
	}
}

func TestDiskRequest(t *testing.T) {
	s := inittests(t)
	disk := s.DiskRequest()
	var expected int64
	expected = 0
	if disk != expected {
		t.Errorf("Disk request was %d, not %d", disk, expected)
	}
}

func TestExtractJobID(t *testing.T) {
	testData := []byte(`1000 job(s) submitted to cluster 100000000.0000.`)
	actual := ExtractJobID(testData)
	expected := []byte("100000000")
	if !bytes.Equal(actual, expected) {
		t.Errorf("extractJobID found %s instead of %s", actual, expected)
	}

	testData = []byte(`asdfadsfadsfadsfa1000 job(s) submitted to cluster 100000000.0000asdfadsfadsfasdfadsfadsfadsfadsfadsf`)
	actual = ExtractJobID(testData)
	expected = []byte("100000000")
	if !bytes.Equal(actual, expected) {
		t.Errorf("extractJobID found %s instead of %s", actual, expected)
	}

	testData = []byte(`asdfadsfadsfadsfa
adsfadsfadsfadsfadsfasdfadsfadsfadsfadsfdsa1000 job(s) submitted to cluster 100000000asdfadsfadsfasdfadsfadsfadsfadsfadsf
asdfadsfasdfadsfdsfsdsfdsafds`)
	actual = ExtractJobID(testData)
	expected = []byte("100000000")
	if !bytes.Equal(actual, expected) {
		t.Errorf("extractJobID found %s instead of %s", actual, expected)
	}

}
