// DES1.cpp: implementation of the CDES class.
//
//////////////////////////////////////////////////////////////////////
#include "common.h"
#include "des.h"
#include "d3des.h"

#ifdef _DEBUG
#undef THIS_FILE
static char THIS_FILE[]=__FILE__;
#endif

//////////////////////////////////////////////////////////////////////
// Construction/Destruction
//////////////////////////////////////////////////////////////////////


bool enc(const char* pass, unsigned char *inbuf, int inlen, unsigned char *outbuf)
{
	const int head_len = sizeof(WORD);
	unsigned char passwd[8] = {0};
	StrToKey(pass, (char*)passwd);
	//	if (strlen(pass) > 8)
	//		memcpy(passwd, pass, 8);
	//	else
	//		memcpy(passwd, pass, strlen(pass));

	int step = inlen / 8;

	deskey(passwd, EN0);
	for (int i=0; i<step; i++)
	{
		des(&inbuf[8*i], &outbuf[8*i + head_len]);
	}

	int r = inlen % 8;
	if (r != 0)
	{
		unsigned char buf[8] = {0};
		memcpy(buf, &inbuf[inlen - r], r);
		des(buf, &outbuf[step*8 + head_len]);
	}

	WORD size = inlen;
	memcpy(outbuf, &size, head_len);
	return true;
}

bool dec(const char* pass, unsigned char *inbuf, int inlen, unsigned char *outbuf, int& outlen)
{
	unsigned char passwd[8] = {0};
	StrToKey(pass, (char*)passwd);
	//	if (strlen(pass) > 8)
	//		memcpy(passwd, pass, 8);
	//	else
	//		memcpy(passwd, pass, strlen(pass));

	//1. �ж�inlen�ĳ���
	if (inlen % 8 != 2)
	{
		return false;
	}

	WORD size;
	memcpy(&size, inbuf, sizeof(WORD));
	ASSERT(size + 2 <= inlen);
	outlen = size;

	int step = (inlen - 2) / 8;

	deskey(passwd, DE1);
	for (int i=0; i<step; i++)
	{
		des(&inbuf[2+8*i], &outbuf[8*i]);
	}
	return true;
}

//��ԿΪ8��byte
bool enc(const unsigned char pass[], unsigned char* inbuf, int inlen, unsigned char* outbuf)
{
	const int head_len = sizeof(WORD);
	unsigned char passwd[8] = {0};
	//StrToKey(pass, (char*)passwd);
	memcpy(passwd, pass, 8);

	int step = inlen / 8;

	deskey(passwd, EN0);
	for (int i=0; i<step; i++)
	{
		des(&inbuf[8*i], &outbuf[8*i + head_len]);
	}

	int r = inlen % 8;
	if (r != 0)
	{
		BYTE buf[8] = {0};
		memcpy(buf, &inbuf[inlen - r], r);
		des(buf, &outbuf[step*8 + head_len]);
	}

	WORD size = inlen;
	memcpy(outbuf, &size, head_len);
	return true;
}

bool dec(const unsigned char pass[], unsigned char* inbuf, int inlen, unsigned char* outbuf, int& outlen)
{
	unsigned char passwd[8] = {0};
	memcpy(passwd, pass, 8);
	//StrToKey(pass, (char*)passwd);

	//1. �ж�inlen�ĳ���
	if (inlen % 8 != 2)
	{
		return false;
	}

	WORD size;
	memcpy(&size, inbuf, sizeof(WORD));
	ASSERT(size + 2 <= inlen);
	outlen = size;

	int step = (inlen - 2) / 8;

	deskey(passwd, DE1);
	for (int i=0; i<step; i++)
	{
		des(&inbuf[2+8*i], &outbuf[8*i]);
	}
	return true;
}

//�������루���Ȳ�����16��
bool enc_passwd(const char* pass, unsigned char data[])
{
	unsigned char passwd[8] = {0};
	StrToKey(pass, (char*)passwd);
	//memcpy(passwd, pass, 8);

	deskey(passwd, EN0);
	des(data, data);
	des(&data[8], &data[8]);
	return true;
}

//�������루���Ȳ�����16��
bool dec_passwd(const char* pass, unsigned char data[])
{
	unsigned char passwd[8] = {0};
	StrToKey(pass, (char*)passwd);
	//	memcpy(passwd, pass, 8);

	deskey(passwd, DE1);
	des(data, data);
	des(&data[8], &data[8]);
	return true;
}


int enc_len(int len)
{
	if (len % 8 == 0)
		return len + 2;
	else
		return (len / 8 + 1)*8 + 2;
}

void StrToKey(const char *str, char *key)
{
	char tmpkey[8];
	int i, k, j;
	int len = strlen(str);
	int count = len / 8;
	int mod = len % 8;

	for (i=0; i<count; i++ )
	{
		for(k=0;k<8;k++)
		{
			tmpkey[k]= str[8*i+k];
			key[k]^=tmpkey[k];
		}
	}

	if (mod != 0)
	{
		for(k=0;k<mod;k++)
		{
			tmpkey[k]=str[len-mod+k];
			key[k]^=tmpkey[k];
		}
		for(j=0;j<8-mod;j++)
		{
			tmpkey[mod+j]=0;
			key[k]^=tmpkey[k];
		}
	}
}
