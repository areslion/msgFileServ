#include <Ice/Ice.h>

#include <fstream>
#include <iostream>

#include "ICEOrg.h"


using namespace std;
using namespace TRMSOrganization;
const int chunksize =1024;




int
main(int argc, char* argv[])
{
	int status = 0;
	Ice::CommunicatorPtr ic;
	
	
	try {
		ic = Ice::initialize(argc, argv,"./config.client");
		Ice::ObjectPrx base = ic->propertyToProxy("OrgEmp.Proxy");
		OrganizationServicePrx vmp = OrganizationServicePrx::checkedCast(base);


		if (!vmp)
			throw "Invalid VersionManager proxy";
	}

#if(0)
		//请求最新版本
		try
		{
			vmp->getLastVersionNumber(v);
		}catch(updateError &ex)
		{
			status = -1;
			cerr << ex.msg << endl;
		}
		if(status )
		{
			cout<<v.major<<"."<<v.minor<<"."<<v.revision<<endl;
		}
		//上传升级文件
		UpdateFile file={
			"final-version.zip",
			{1,1,11},
			"ae2ab64d9972168813528dd850ce8a8d",
			11240
		};
		ifile.open(file.fileName.c_str(), ios::in | ios::binary);
		if (!ifile.is_open())
		{
			 throw string("Can't open file "+file.fileName+"!");
		}
		ifile.seekg(0,ios::end);
		file.FileSize = ifile.tellg();
		ifile.seekg(0, ios::beg);
		cout<< "file size is " <<file.FileSize<<endl;
		try
		{							   
			vmp->requestUploadUpdateFile(file);
			while (!ifile.eof())
			{
				ifile.read(buf, chunksize);
				readSize=ifile.gcount();
				
				bs.reserve(readSize);
				bs.assign(&buf[0], &buf[readSize]);
				//发送文件
				vmp->sendFile(bs);
				
				sendLen+=readSize;
				
				percentage = (int)(((float)sendLen/file.FileSize)*100);
				
				cout << "send file ..."<<percentage<<"%\r";
				cout.flush();
			}
			cout<<endl;
			ifile.close();
			cout << "done!"<< endl;
		}catch(updateError &ex)
		{
			//cerr << ex.msg << endl;
			throw ex.msg;
		}
		catch(fileTransferError &ex)
		{
			throw ex.msg;
		}
		
	} 
	catch (const Ice::Exception& ex) {
		cerr << ex << endl;
		status = 1;
	}
	catch (const char* msg) {
		cerr << msg << endl;
		status = 1;
	}
	catch (const string msg)
	{
		cerr<< msg <<endl;
		status = 1;
	}
	cout << "done2!"<< endl;
	if (ic)
	ic->destroy();
	ifile.close();
	cout << "done3!"<< endl;
#endif
	catch (const Ice::Exception& ex) {
		cerr << ex << endl;
		status = 1;
	}
	catch (const string msg)
       {
		cerr<< msg <<endl;

		status = 1;
	}

	sleep(100);

	return status;
}
