/*
SQLyog Ultimate v12.09 (64 bit)
MySQL - 5.5.56-MariaDB : Database - deskSafe
*********************************************************************
*/

/*!40101 SET NAMES utf8 */;

/*!40101 SET SQL_MODE=''*/;

/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
CREATE DATABASE /*!32312 IF NOT EXISTS*/`deskSafe` /*!40100 DEFAULT CHARACTER SET utf8 */;

USE `deskSafe`;

/*Table structure for table `depotSft` */

DROP TABLE IF EXISTS `depotSft`;

CREATE TABLE `depotSft` (
  `namexf` char(60) NOT NULL COMMENT '软件名称 不能重复',
  `namexa` char(60) NOT NULL COMMENT '软件显示名称',
  `ver` char(20) NOT NULL COMMENT '软件版本',
  `pathx` text NOT NULL COMMENT '安装包资源路径',
  `pathIcon` text COMMENT 'Icon图标路径',
  `flagSft` int(4) unsigned NOT NULL COMMENT '0x01办公 0x02常用 0x04必备',
  `md5x` char(40) DEFAULT NULL COMMENT '软件包的md5',
  `folderId` char(40) NOT NULL COMMENT '软件资源所在的文件夹 guid命名',
  `descx` text COMMENT '安装软件的描述信息',
  PRIMARY KEY (`namexa`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

/*Table structure for table `devMoveableInfo` */

DROP TABLE IF EXISTS `devMoveableInfo`;

CREATE TABLE `devMoveableInfo` (
  `numEvt` char(40) CHARACTER SET gb2312 NOT NULL COMMENT '产生事件的ID',
  `numDev` char(40) CHARACTER SET gb2312 NOT NULL COMMENT '终端电脑ID',
  `numHdev` char(40) CHARACTER SET gb2312 NOT NULL COMMENT '终端设备ID',
  `numUser` char(40) NOT NULL COMMENT '用户UKEY',
  `tmEvt` datetime NOT NULL COMMENT '事件产生的时间',
  `devName` char(40) CHARACTER SET gb2312 NOT NULL COMMENT '设备名称',
  `devCls` int(4) unsigned NOT NULL COMMENT '设备类型 ukey unknow',
  `devDesc` char(60) CHARACTER SET gb2312 NOT NULL COMMENT '设备描述信息',
  `EvtCls` int(4) NOT NULL COMMENT '事件类型 0 为止 1 插入 2拔出',
  `ipr` char(40) NOT NULL COMMENT '远端客户端IP地址',
  `flagx` int(4) unsigned NOT NULL DEFAULT '0',
  `rev1` char(40) CHARACTER SET gb2312 DEFAULT NULL,
  `rev2` char(40) CHARACTER SET gb2312 DEFAULT NULL,
  `rev3` char(40) CHARACTER SET gb2312 DEFAULT NULL,
  PRIMARY KEY (`numEvt`),
  KEY `tmEvt` (`tmEvt`,`flagx`),
  KEY `flagx` (`flagx`),
  KEY `numDev` (`numDev`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*Table structure for table `employee` */

DROP TABLE IF EXISTS `employee`;

CREATE TABLE `employee` (
  `numSelf` char(40) CHARACTER SET gb2312 NOT NULL COMMENT '人员编号',
  `numOrgnization` char(40) CHARACTER SET gb2312 NOT NULL COMMENT '所在的机构编码',
  `numPosition` char(40) CHARACTER SET gb2312 NOT NULL COMMENT '职位编码',
  `numLeader` char(40) CHARACTER SET gb2312 NOT NULL COMMENT '主管领导编码',
  `numCell` char(16) NOT NULL DEFAULT '""' COMMENT '电话号码',
  `ukey` char(40) CHARACTER SET gb2312 NOT NULL COMMENT '绑定的人员终端的ID',
  `namex` char(80) NOT NULL COMMENT '职员姓名',
  `email` char(60) CHARACTER SET gb2312 NOT NULL COMMENT '邮箱',
  `pwdlogin` char(60) NOT NULL COMMENT '登陆密码',
  `gender` int(4) NOT NULL COMMENT '性别,1为男性，0为女性',
  `tmCreate` datetime NOT NULL COMMENT '创建时间',
  `tmMod` datetime NOT NULL COMMENT '修改时间',
  `priviege` int(12) unsigned NOT NULL DEFAULT '1' COMMENT '权限',
  `pathx` varchar(1024) DEFAULT NULL COMMENT '所属部门路径',
  `rever1` char(40) CHARACTER SET gb2312 DEFAULT NULL,
  `rever2` char(40) CHARACTER SET gb2312 DEFAULT NULL,
  `rever3` char(40) CHARACTER SET gb2312 DEFAULT NULL,
  PRIMARY KEY (`numCell`,`email`),
  KEY `numCell` (`numCell`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*Table structure for table `employeeGroup` */

DROP TABLE IF EXISTS `employeeGroup`;

CREATE TABLE `employeeGroup` (
  `namex` char(80) NOT NULL COMMENT '职员姓名',
  `ukey` char(40) CHARACTER SET gb2312 NOT NULL COMMENT '绑定的人员终端的ID',
  `email` char(60) CHARACTER SET gb2312 NOT NULL COMMENT '邮箱',
  `gender` int(4) NOT NULL COMMENT '性别,1为男性，0为女性',
  `tmCreate` datetime NOT NULL COMMENT '创建时间',
  `tmMod` datetime NOT NULL COMMENT '修改时间',
  `priviege` int(12) unsigned NOT NULL DEFAULT '1' COMMENT '权限',
  `pathx` varchar(1024) DEFAULT NULL COMMENT '所属部门路径',
  `rever1` char(40) CHARACTER SET gb2312 DEFAULT NULL,
  `rever2` char(40) CHARACTER SET gb2312 DEFAULT NULL,
  `rever3` char(40) CHARACTER SET gb2312 DEFAULT NULL,
  PRIMARY KEY (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*Table structure for table `hotFix` */

DROP TABLE IF EXISTS `hotFix`;

CREATE TABLE `hotFix` (
  `numEvt` char(40) NOT NULL COMMENT '事件ID',
  `numDev` char(40) NOT NULL COMMENT '终端设备ID',
  `numHotFix` char(40) NOT NULL COMMENT '补丁ID',
  `namex` char(120) NOT NULL COMMENT '补丁名称',
  `description` varchar(512) NOT NULL COMMENT '补丁描述信息',
  `tmEvt` datetime NOT NULL COMMENT '报表事件时间',
  `tmFix` datetime NOT NULL COMMENT '补丁安装时间',
  `tmGet` datetime DEFAULT NULL COMMENT '采集时间',
  `flagx` int(4) unsigned NOT NULL DEFAULT '0',
  `rev1` char(40) DEFAULT NULL,
  `rev2` char(40) DEFAULT NULL,
  `rev3` char(40) DEFAULT NULL,
  PRIMARY KEY (`numEvt`),
  KEY `tmEvt` (`tmEvt`),
  KEY `flagx` (`flagx`),
  KEY `numDev` (`numDev`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

/*Table structure for table `msgAbstract` */

DROP TABLE IF EXISTS `msgAbstract`;

CREATE TABLE `msgAbstract` (
  `numMsg` char(40) NOT NULL COMMENT '消息GUID',
  `namex` char(120) NOT NULL COMMENT '消息名称',
  `tmx` datetime NOT NULL COMMENT '执行开始时间',
  `tmy` datetime NOT NULL COMMENT '执行结束时间',
  `tmm` datetime NOT NULL COMMENT '创建修改时间',
  `os` int(4) unsigned NOT NULL COMMENT '操作系统类型',
  `autoexe` tinyint(1) NOT NULL DEFAULT '0' COMMENT '自动执行',
  `popup` tinyint(1) NOT NULL DEFAULT '0' COMMENT '弹窗通知',
  `numSender` char(40) NOT NULL COMMENT '发送人key',
  `numSent` int(4) unsigned NOT NULL DEFAULT '0' COMMENT '分发次数',
  `numSentOK` int(4) unsigned NOT NULL DEFAULT '0' COMMENT '发送成功次数',
  `numSentKO` int(4) unsigned NOT NULL DEFAULT '0' COMMENT '发送失败数',
  `descx` text COMMENT '任务描述信息',
  PRIMARY KEY (`numMsg`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

/*Table structure for table `msgSend` */

DROP TABLE IF EXISTS `msgSend`;

CREATE TABLE `msgSend` (
  `numMsg` char(40) NOT NULL COMMENT '消息GUID',
  `numReciever` char(40) NOT NULL COMMENT '接收者KEY',
  `statusx` int(4) NOT NULL DEFAULT '1' COMMENT '0执行成功 1客户端收到 2客户端执行成功 3客户端执行失败 99所有状态',
  `os` int(4) NOT NULL DEFAULT '0' COMMENT '操作系统',
  `tmm` datetime NOT NULL COMMENT '消息产生时间',
  `tmSend` datetime DEFAULT NULL COMMENT '消息发送时间',
  `tmExc` datetime DEFAULT NULL COMMENT '消息执行成功时间',
  `tmx` datetime DEFAULT NULL COMMENT '消息有效开始时间',
  `tmy` datetime DEFAULT NULL COMMENT '消息有效结束时间',
  `descx` text COMMENT '消息描述信息',
  `namex` char(120) DEFAULT NULL COMMENT '消息名称',
  PRIMARY KEY (`numMsg`,`numReciever`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

/*Table structure for table `oprand` */

DROP TABLE IF EXISTS `oprand`;

CREATE TABLE `oprand` (
  `numEvt` char(40) CHARACTER SET gb2312 NOT NULL COMMENT '数据库产生的GUID',
  `tmxC` datetime NOT NULL COMMENT '采集时间',
  `devID` char(120) CHARACTER SET gb2312 NOT NULL COMMENT '设备的ID',
  `oprandx` char(60) CHARACTER SET gb2312 NOT NULL COMMENT '操作码字符创表示',
  `tmx` int(8) unsigned NOT NULL COMMENT '操作开始时间',
  `tmy` int(8) unsigned NOT NULL COMMENT '操作结束时间',
  `tmspan` int(8) unsigned NOT NULL COMMENT '操作时长微毫为单位',
  `rever1` char(40) CHARACTER SET gb2312 DEFAULT NULL,
  `rever2` char(40) CHARACTER SET gb2312 DEFAULT NULL,
  `rever3` char(40) CHARACTER SET gb2312 DEFAULT NULL,
  PRIMARY KEY (`numEvt`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

/*Table structure for table `orgnization` */

DROP TABLE IF EXISTS `orgnization`;

CREATE TABLE `orgnization` (
  `numSelf` char(60) CHARACTER SET gb2312 NOT NULL COMMENT '本级机构编码',
  `numFather` char(60) CHARACTER SET gb2312 NOT NULL COMMENT '本级机构父节点编码',
  `name` text CHARACTER SET gb2312 NOT NULL COMMENT '本级机构名称',
  `level` int(4) unsigned NOT NULL COMMENT '层级0~n,0为最高层级',
  `tmx` datetime NOT NULL COMMENT '创建时间',
  `tmy` datetime NOT NULL COMMENT '修改时间',
  `address` text CHARACTER SET gb2312 COMMENT '机构所在地',
  `description` text CHARACTER SET gb2312 COMMENT '描述信息',
  `rever1` char(40) CHARACTER SET gb2312 DEFAULT NULL,
  `rever2` char(40) CHARACTER SET gb2312 DEFAULT NULL,
  `rever3` char(40) CHARACTER SET gb2312 DEFAULT NULL,
  PRIMARY KEY (`numSelf`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*Table structure for table `portx` */

DROP TABLE IF EXISTS `portx`;

CREATE TABLE `portx` (
  `numEvt` char(40) CHARACTER SET gb2312 NOT NULL COMMENT '事件ID',
  `numDev` char(40) CHARACTER SET gb2312 NOT NULL COMMENT '终端ID',
  `tmEvt` datetime NOT NULL COMMENT '报表事件时间',
  `tmGet` datetime NOT NULL COMMENT '采集时间',
  `numPort` char(16) CHARACTER SET gb2312 NOT NULL COMMENT '端口ID',
  `portCls` char(16) CHARACTER SET gb2312 NOT NULL COMMENT '端口类型,udp ,tcp',
  `pathx` varchar(512) CHARACTER SET gb2312 NOT NULL COMMENT '端口事件进程路径',
  `ipr` char(40) NOT NULL,
  `flagx` int(4) unsigned NOT NULL DEFAULT '0',
  `rev1` char(40) CHARACTER SET gb2312 DEFAULT NULL,
  `rev2` char(40) CHARACTER SET gb2312 DEFAULT NULL,
  `rev3` char(40) CHARACTER SET gb2312 DEFAULT NULL,
  PRIMARY KEY (`numEvt`),
  KEY `tmEvt` (`tmEvt`),
  KEY `flagx` (`flagx`),
  KEY `numDev` (`numDev`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*Table structure for table `positionLevel` */

DROP TABLE IF EXISTS `positionLevel`;

CREATE TABLE `positionLevel` (
  `numPosition` char(40) CHARACTER SET gb2312 NOT NULL COMMENT '职位编码',
  `numPostFater` char(40) CHARACTER SET gb2312 NOT NULL COMMENT '所属上级职位编码',
  `description` text CHARACTER SET gb2312 COMMENT '职位描述',
  `privilege` int(8) unsigned NOT NULL COMMENT '职位权限',
  `rever1` char(40) CHARACTER SET gb2312 DEFAULT NULL,
  `rever2` char(40) CHARACTER SET gb2312 DEFAULT NULL,
  `rever3` char(40) CHARACTER SET gb2312 DEFAULT NULL,
  PRIMARY KEY (`numPosition`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*Table structure for table `processx` */

DROP TABLE IF EXISTS `processx`;

CREATE TABLE `processx` (
  `numEvt` char(40) CHARACTER SET gb2312 NOT NULL COMMENT '采集事件的ID',
  `numDev` char(40) CHARACTER SET gb2312 NOT NULL COMMENT '设备ID',
  `numProc` char(16) CHARACTER SET gb2312 NOT NULL COMMENT '进程ID',
  `nameProc` char(40) CHARACTER SET gb2312 NOT NULL COMMENT '进程名称',
  `tmEvt` datetime NOT NULL COMMENT '报表事件时间',
  `tmGet` datetime NOT NULL COMMENT '采集时间',
  `pathx` varchar(512) CHARACTER SET gb2312 NOT NULL COMMENT '进程路径',
  `useageCPU` char(8) CHARACTER SET gb2312 NOT NULL COMMENT 'CPU使用率',
  `useageMem` char(8) CHARACTER SET gb2312 NOT NULL COMMENT '内存使用率',
  `flagx` int(4) unsigned NOT NULL DEFAULT '0',
  `rev1` char(40) CHARACTER SET gb2312 DEFAULT NULL,
  `rev2` char(40) CHARACTER SET gb2312 DEFAULT NULL,
  `rev3` char(40) CHARACTER SET gb2312 DEFAULT NULL,
  PRIMARY KEY (`numEvt`),
  KEY `tmEvt` (`tmEvt`),
  KEY `flagx` (`flagx`),
  KEY `numDev` (`numDev`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*Table structure for table `property` */

DROP TABLE IF EXISTS `property`;

CREATE TABLE `property` (
  `numEvt` char(40) NOT NULL COMMENT '事件ID',
  `numDev` char(40) NOT NULL COMMENT '资产ID(终端设备ID)',
  `numUsr` char(40) NOT NULL COMMENT '资产归属人的ID',
  `numUsrBy` char(40) NOT NULL COMMENT '资产变更操作人ID 为空的时候是自动关联',
  `flagx` int(4) unsigned NOT NULL DEFAULT '0' COMMENT '标记 0 资产正常  1  资产停用',
  `tmEvt` datetime NOT NULL COMMENT '事件产生的时间',
  PRIMARY KEY (`numEvt`),
  KEY `numDev` (`numDev`),
  KEY `numUsr` (`numUsr`),
  KEY `flagx` (`flagx`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*Table structure for table `sftRuleAbstract` */

DROP TABLE IF EXISTS `sftRuleAbstract`;

CREATE TABLE `sftRuleAbstract` (
  `num` char(40) NOT NULL COMMENT '策略GUID',
  `namex` char(120) NOT NULL COMMENT '策略名称',
  `cls` int(4) NOT NULL COMMENT '策略分类',
  PRIMARY KEY (`num`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

/*Table structure for table `sftRuleSend` */

DROP TABLE IF EXISTS `sftRuleSend`;

CREATE TABLE `sftRuleSend` (
  `numRule` char(40) COLLATE utf8_estonian_ci NOT NULL COMMENT '策略的GUID',
  `numUser` char(40) COLLATE utf8_estonian_ci NOT NULL COMMENT '分发对象的GUID',
  `namex` char(120) COLLATE utf8_estonian_ci NOT NULL COMMENT '策略名称',
  `cls` int(4) NOT NULL COMMENT '策略分类',
  `path` text COLLATE utf8_estonian_ci COMMENT '所在部门',
  PRIMARY KEY (`numRule`,`numUser`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_estonian_ci;

/*Table structure for table `sharedDirs` */

DROP TABLE IF EXISTS `sharedDirs`;

CREATE TABLE `sharedDirs` (
  `numEvt` char(40) NOT NULL COMMENT '采集事件的ID',
  `numDev` char(40) NOT NULL COMMENT '终端设备ID',
  `tmEvt` datetime NOT NULL COMMENT '报表事件时间',
  `tmGet` datetime NOT NULL COMMENT '信息采集时间',
  `ipr` char(40) NOT NULL COMMENT '终端IP',
  `pathx` varchar(512) NOT NULL COMMENT '共享路径',
  `notes` varchar(512) DEFAULT NULL COMMENT '备注信息',
  `flagx` int(4) unsigned NOT NULL DEFAULT '0',
  `rev1` char(40) DEFAULT NULL,
  `rev2` char(40) DEFAULT NULL,
  `rev3` char(40) DEFAULT NULL,
  PRIMARY KEY (`numEvt`),
  KEY `tmEvt` (`tmEvt`),
  KEY `flagx` (`flagx`),
  KEY `numDev` (`numDev`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

/*Table structure for table `softInstalled` */

DROP TABLE IF EXISTS `softInstalled`;

CREATE TABLE `softInstalled` (
  `numEvt` char(40) NOT NULL COMMENT '事件ID',
  `numDev` char(40) NOT NULL COMMENT '终端ID',
  `tmEvt` datetime NOT NULL COMMENT '报表时间',
  `tmGet` datetime DEFAULT NULL COMMENT '获取时间',
  `namex` varchar(512) NOT NULL COMMENT '软件名称',
  `publisher` varchar(512) NOT NULL COMMENT '发布者',
  `ipr` char(32) DEFAULT NULL COMMENT '远端终端IP',
  `ver` char(32) NOT NULL COMMENT '版本',
  `siezx` char(32) NOT NULL COMMENT '软件大小',
  `tmInstall` datetime NOT NULL COMMENT '安装日期',
  `pathx` varchar(512) DEFAULT NULL COMMENT '安装路径',
  `flagx` int(4) unsigned NOT NULL DEFAULT '0',
  `rev1` char(40) DEFAULT NULL,
  `rev2` char(40) DEFAULT NULL,
  `rev3` char(40) DEFAULT NULL,
  PRIMARY KEY (`numEvt`),
  KEY `tmEvt` (`tmEvt`),
  KEY `flagx` (`flagx`),
  KEY `numDev` (`numDev`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

/*Table structure for table `terDevBasicInfo` */

DROP TABLE IF EXISTS `terDevBasicInfo`;

CREATE TABLE `terDevBasicInfo` (
  `numDev` char(40) CHARACTER SET gb2312 NOT NULL COMMENT '设备的ID',
  `numUsrKey` char(40) CHARACTER SET gb2312 NOT NULL DEFAULT 'NULL' COMMENT '终端责任人的UKEY',
  `numPwdHelp` char(10) NOT NULL COMMENT '远程协助密码',
  `tmCreate` datetime NOT NULL COMMENT '第一次录入时间',
  `tmMod` datetime NOT NULL COMMENT '设备信息修改时间(不含心跳)',
  `netCardAddr` char(40) CHARACTER SET gb2312 NOT NULL COMMENT '有线物理网卡地址',
  `netCardInfo` char(120) CHARACTER SET gb2312 NOT NULL COMMENT '网卡信息',
  `verA` int(8) NOT NULL COMMENT '大版本号',
  `verB` int(8) NOT NULL COMMENT '中版本号',
  `verC` int(8) NOT NULL COMMENT '小版本号',
  `cpuCls` char(40) CHARACTER SET gb2312 NOT NULL COMMENT 'CPU类型',
  `cpusName` char(40) CHARACTER SET gb2312 NOT NULL COMMENT 'CPU名字',
  `diskHName` char(40) CHARACTER SET gb2312 NOT NULL COMMENT '硬盘名称',
  `diskCls` char(40) CHARACTER SET gb2312 NOT NULL COMMENT '硬盘类型',
  `diskSize` int(4) NOT NULL COMMENT '硬盘大小，G为单位',
  `mainBoardM` char(40) CHARACTER SET gb2312 NOT NULL COMMENT '主板厂商',
  `devName` char(40) CHARACTER SET gb2312 NOT NULL COMMENT '计算机名字',
  `devStatus` int(4) unsigned NOT NULL COMMENT '设备状态(0停用 1启用)',
  `systeminfo` char(40) CHARACTER SET gb2312 NOT NULL COMMENT '操作系统描述',
  `ipx` char(40) NOT NULL COMMENT 'IP地址',
  `flagx` int(4) unsigned NOT NULL DEFAULT '0',
  `bEnableHelp` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否允许远程协助，默认为否',
  `rev1` char(40) CHARACTER SET gb2312 DEFAULT NULL,
  `rev2` char(40) CHARACTER SET gb2312 DEFAULT NULL,
  `rev3` char(40) CHARACTER SET gb2312 DEFAULT NULL,
  PRIMARY KEY (`numDev`),
  KEY `flagx` (`flagx`),
  KEY `numDev` (`numDev`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

/*Table structure for table `terminalDevRunInfo` */

DROP TABLE IF EXISTS `terminalDevRunInfo`;

CREATE TABLE `terminalDevRunInfo` (
  `numEvt` char(40) CHARACTER SET gb2312 NOT NULL,
  `numDev` char(40) CHARACTER SET gb2312 NOT NULL,
  `numUsrKey` char(40) DEFAULT NULL COMMENT '用户插入的UKEY信息',
  `ipr` char(32) DEFAULT NULL COMMENT '远端客户端的自己获取的IP地址',
  `ipFrmGate` char(32) NOT NULL COMMENT '从网关解析出来的IP地址',
  `tmEvt` datetime NOT NULL COMMENT '报表事件时间',
  `tmGet` datetime NOT NULL COMMENT '采集时间',
  `hdiskTotal` int(4) unsigned NOT NULL,
  `hdiskUsed` int(4) unsigned NOT NULL,
  `hdiskUseage` float unsigned NOT NULL,
  `cpuUseage` char(60) CHARACTER SET gb2312 NOT NULL,
  `memTotal` int(4) unsigned NOT NULL,
  `memUsed` int(4) unsigned NOT NULL,
  `memUseage` float unsigned DEFAULT NULL,
  `netSpeedUP` int(4) unsigned NOT NULL,
  `netSpeedDown` int(4) unsigned NOT NULL,
  `verTerm` char(30) NOT NULL COMMENT '企管通客户端版本号',
  `flagx` int(4) unsigned NOT NULL DEFAULT '0',
  `rever1` char(40) CHARACTER SET gb2312 DEFAULT NULL,
  `rever2` char(40) CHARACTER SET gb2312 DEFAULT NULL,
  `rever3` char(40) CHARACTER SET gb2312 DEFAULT NULL,
  PRIMARY KEY (`numEvt`),
  KEY `tmEvt` (`tmEvt`),
  KEY `flagx` (`flagx`),
  KEY `numDev` (`numDev`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

/*Table structure for table `windowApp` */

DROP TABLE IF EXISTS `windowApp`;

CREATE TABLE `windowApp` (
  `numEvt` char(40) CHARACTER SET gb2312 NOT NULL COMMENT '事件ID',
  `numDev` char(40) CHARACTER SET gb2312 NOT NULL COMMENT '设备ID',
  `tmEvt` datetime NOT NULL COMMENT '报表事件时间',
  `tmGet` datetime NOT NULL COMMENT '采集时间',
  `title` text CHARACTER SET gb2312 NOT NULL COMMENT '窗口标题',
  `flagx` int(4) unsigned NOT NULL DEFAULT '0',
  `rev1` char(40) CHARACTER SET gb2312 DEFAULT NULL,
  `rev2` char(40) CHARACTER SET gb2312 DEFAULT NULL,
  `rev3` char(40) CHARACTER SET gb2312 DEFAULT NULL,
  PRIMARY KEY (`numEvt`),
  KEY `tmEvt` (`tmEvt`),
  KEY `flagx` (`flagx`),
  KEY `numDev` (`numDev`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
