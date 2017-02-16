package update

import (
	"time"
	"math/rand"
	"path/filepath"
	"os"
	"strings"
	"fmt"
)


/*

  #define EXEC_TIMEOUT 60	//execute the command timeout 60s
  #define UPD_TIMEOUT	1800//execute the update timeout 30 minutes
  #define CONN_TIMEOUT 120	//connect time out be 2 minutes
  EXEC_TIMEOUT = 60
  UPD_TIMEOUT = 1800
  CONN_TIMEOUT = 120
  APPVERSION_FILE = "/app/appversion"
  CATMAC = "cat /usr/sbin/macaddr"
  #define  APP_VERSION_FILE "/app/appversion"
  #define  CFG_VERSION_FILE "/config/cfgversion"
  #define  SVR_APPPRE_FILE	 "/etc/dlancmd/apppre"
  #define  SVR_APPPRE_FILE_ARM "/var/dlancmd/apppre"
  #define  SVR_APPPKG_FILE "/stmp/app"
  #define  SVR_APPSH_FILE		 "/etc/dlancmd/appsh"
  #define  SVR_APPSH_FILE_ARM  "/var/dlancmd/appsh"
  #define  SVR_CFGPRE_FILE	 "/etc/dlancmd/cfgpre"
  #define  SVR_CFGPRE_FILE_ARM "/var/dlancmd/cfgpre"
  #define  SVR_CFGPKG_FILE "/stmp/cfg"
  #define  SVR_CFGSH_FILE		 "/etc/dlancmd/cfgsh"
  #define  SVR_CFGSH_FILE_ARM  "/var/dlancmd/cfgsh"
  #define  SVR_RECOVPRE_FILE	 "/etc/dlancmd/prercovcfgsh"
  #define  SVR_RECOVPRE_FILE_ARM "/var/dlancmd/prercovcfgsh"
  #define  NEW_SVR_RECOVPRE_FILE	 "/usr/sbin/prercovcfgsh"
  #define  SVR_RECOVPKG_FILE "/stmp/cfgbk"
  #define  SVR_RECOVSH_FILE		"/etc/dlancmd/rcovcfgsh"
  #define  SVR_RECOVSH_FILE_ARM  "/var/dlancmd/rcovcfgsh"
  #define  NEW_SVR_RECOVSH_FILE	"/usr/sbin/rcovcfgsh"
  #define  SVR_BAKPKG_FILE "/stmp/cfgbk"
  #define  SVR_BAKSH_FILE		 "/etc/dlancmd/bakcfgsh"
  #define  SVR_BAKSH_FILE_ARM  "/var/dlancmd/bakcfgsh"
  #define  NEW_SVR_BAKSH_FILE	 "/usr/sbin/bakcfgsh"
  #define  SVR_RESULT_FILE		"/etc/dlancmd/result"
  #define  SVR_RESULT_FILE_ARM	"/var/dlancmd/result"
  #define SVR_RESULT_FILE_TIME		"/etc/dlancmd/result_time"	//add 6.0
  #define SVR_RESULT_FILE_TIME_ARM	"/var/dlancmd/result_time"	//add 6.0
  #define  SVR_UPDHISTORY_FILE		"/usr/sbin/updhistory.sh"
  #define  SVR_UPDHISTORY_FILE_ARM	"/sbin/updhistory.sh"
  #define  SVR_SHELL_FILE			"/etc/dlancmd/compose.sh" //升级组合包时的清理脚本都上传成compose.sh
  #define  SVR_SHELL_FILE_ARM		"/var/dlancmd/compose.sh"

  #define SVR_SHELL_DNLD			"/etc/dlancmd/dnldsh.sh"	//add 6.0
  #define SVR_SHELL_DNLD_ARM		"/var/dlancmd/dnldsh.sh"

  ARM_LINUX_BASIC = ["/var/dlancmd/tempexec","/var/dlancmd/result","/var/upd_sh_err.log","/var/dlancmd/return","/etc/config/passwd","/var/dlancmd/compose.sh"]
  X86_LINUX_BASIC = ["/etc/dlancmd/tempexec","/etc/dlancmd/result","/var/upd_sh_err.log","/etc/dlancmd/return","/config/passwd","/etc/dlancmd/compose.sh"]
  ARM_LINUX_UPDATE = ["/var/dlancmd/apppre","/var/dlancmd/appsh","/var/dlancmd/cfgpre","/var/dlancmd/cfgsh"]
  X86_LINUX_UPDATE = ["/etc/dlancmd/apppre","/etc/dlancmd/appsh","/etc/dlancmd/cfgpre","/etc/dlancmd/cfgsh"]

  DES_KEY = "dlandproxy"
  SSU_DEC_PASSWD  = "sangforupd~!@#\$%"
  SSU_DEC_PASSWD_OLD  = "greatsinfor"
  CHECK_UPGRADE_SN   = "/app/usr/sbin/checkupdsn.sh"
  CSSU_PACKAGE_CONF = "upgrade.conf"
  SSU_PACKAGE_CONF = "package.conf"
  UPDHISTORY_SCRIPT = "/usr/sbin/updhistory.sh"
  UPDATE_CHECK_SCRIPT = "/usr/sbin/updatercheck.sh"
  BACKUP_SCTRIPT = "/usr/sbin/bakcfgsh"
  PRERECOVCFGSH_SCTRIPT = "/usr/sbin/prercovcfgsh"
  RECOVCFGSH_SCTRIPT = "/usr/sbin/rcovcfgsh"
  PACKAGE_TYPE = 1
  RESTORE_TYPE = 2
  EXECUTE_TYPE = 3
  AUTOBAK_NUMS = 10

 */


var (
	EXEC_TIMEOUT = 60
	UPD_TIMEOUT = 1800
	CONN_TIMEOUT = 120
	APPVERSION_FILE = "/app/appversion"
	CATMAC = "cat /usr/sbin/macaddr"
	ARM_LINUX_BASIC = [6]string{"/var/dlancmd/tempexec","/var/dlancmd/result","/var/upd_sh_err.log","/var/dlancmd/return","/etc/config/passwd","/var/dlancmd/compose.sh" }

	X86_LINUX_BASIC = [6]string{"/etc/dlancmd/tempexec","/etc/dlancmd/result","/var/upd_sh_err.log","/etc/dlancmd/return","/config/passwd","/etc/dlancmd/compose.sh"}
	ARM_LINUX_UPDATE = [4]string{"/var/dlancmd/apppre","/var/dlancmd/appsh","/var/dlancmd/cfgpre","/var/dlancmd/cfgsh"}
	X86_LINUX_UPDATE = [4]string{"/etc/dlancmd/apppre","/etc/dlancmd/appsh","/etc/dlancmd/cfgpre","/etc/dlancmd/cfgsh"}

	DES_KEY = "dlandproxy"
	SSU_DEC_PASSWD  = "sangforupd~!@#$%"
	SSU_DEC_PASSWD_OLD  = "greatsinfor"
	CHECK_UPGRADE_SN   = "/app/usr/sbin/checkupdsn.sh"
	CSSU_PACKAGE_CONF = "upgrade.conf"
	SSU_PACKAGE_CONF = "package.conf"
	UPDHISTORY_SCRIPT = "/usr/sbin/updhistory.sh"
	UPDATE_CHECK_SCRIPT = "/usr/sbin/updatercheck.sh"
	BACKUP_SCTRIPT = "/usr/sbin/bakcfgsh"
	PRERECOVCFGSH_SCTRIPT = "/usr/sbin/prercovcfgsh"
	RECOVCFGSH_SCTRIPT = "/usr/sbin/rcovcfgsh"
	PACKAGE_TYPE = 1
	RESTORE_TYPE = 2
	EXECUTE_TYPE = 3
	AUTOBAK_NUMS = 10
)

//生成随机字符串
func GetRandomString(length int) string{
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("Get Current Directory wrong:",err)
		return "/tmp"
	}
	return strings.Replace(dir, "\\", "/", -1)
}


func FtpDownloadSSUPackage() {

}