package update

/*
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <signal.h>
#include <string.h>
#include <sys/param.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <sys/socket.h>
#include <sys/wait.h>
#include <assert.h>
#include <fcntl.h>
#include <errno.h>

typedef unsigned int    DWORD;
typedef unsigned int    ULONG;
typedef unsigned char   UCHAR;
typedef unsigned short  USHORT;
typedef signed short	WORD;
typedef unsigned int    HANDLE;
typedef unsigned int    UINT;
typedef void*		PVOID;
typedef unsigned char   BOOLEAN;
typedef unsigned char   BOOL;
typedef void		VOID;
typedef	unsigned int*	PULONG;
typedef long		LONG ;
typedef char		CHAR ;
typedef unsigned char	BYTE;

#define ASSERT assert


#define TRUE 1
#define FALSE 0

#define EN0	0	 //MODE == encrypt
#define DE1	1	 //MODE == decrypt

extern void deskey(unsigned char *, int);
extern void usekey(unsigned long *);
extern void cpkey(unsigned long *);
extern void des(unsigned char *, unsigned char *);

static void scrunch(unsigned char *, unsigned long *);
static void unscrun(unsigned long *, unsigned char *);
static void desfunc(unsigned long *, unsigned long *);
static void cookey(unsigned long *);

static unsigned long KnL[32] = { 0L };

static unsigned short bytebit[8]	= {
	01, 02, 04, 010, 020, 040, 0100, 0200 };

static unsigned long bigbyte[24] = {
	0x800000L,	0x400000L,	0x200000L,	0x100000L,
	0x80000L,	0x40000L,	0x20000L,	0x10000L,
	0x8000L,	0x4000L,	0x2000L,	0x1000L,
	0x800L, 	0x400L, 	0x200L, 	0x100L,
	0x80L,		0x40L,		0x20L,		0x10L,
	0x8L,		0x4L,		0x2L,		0x1L	};

//Use the key schedule specified in the Standard (ANSI X3.92-1981).

static unsigned char pc1[56] = {
56, 48, 40, 32, 24, 16,  8,	 0, 57, 49, 41, 33, 25, 17,
9,  1, 58, 50, 42, 34, 26,	18, 10,  2, 59, 51, 43, 35,
62, 54, 46, 38, 30, 22, 14,	 6, 61, 53, 45, 37, 29, 21,
13,  5, 60, 52, 44, 36, 28,	20, 12,  4, 27, 19, 11,  3 };

static unsigned char totrot[16] = {
1,2,4,6,8,10,12,14,15,17,19,21,23,25,27,28 };

static unsigned char pc2[48] = {
13, 16, 10, 23,  0,  4,  2, 27, 14,  5, 20,  9,
22, 18, 11,  3, 25,  7, 15,  6, 26, 19, 12,  1,
40, 51, 30, 36, 46, 54, 29, 39, 50, 44, 32, 47,
43, 48, 38, 55, 33, 52, 45, 41, 49, 35, 28, 31 };

//Thanks to James Gillogly & Phil Karn!
void deskey(unsigned char *key,int edf)
{
register int i, j, l, m, n;
unsigned char pc1m[56], pcr[56];
unsigned long kn[32];

for ( j = 0; j < 56; j++ ) {
l = pc1[j];
m = l & 07;
pc1m[j] = (key[l >> 3] & bytebit[m]) ? 1 : 0;
}
for( i = 0; i < 16; i++ ) {
if( edf == DE1 ) m = (15 - i) << 1;
else m = i << 1;
n = m + 1;
kn[m] = kn[n] = 0L;
for( j = 0; j < 28; j++ ) {
l = j + totrot[i];
if( l < 28 ) pcr[j] = pc1m[l];
else pcr[j] = pc1m[l - 28];
}
for( j = 28; j < 56; j++ ) {
l = j + totrot[i];
if( l < 56 ) pcr[j] = pc1m[l];
else pcr[j] = pc1m[l - 28];
}
for( j = 0; j < 24; j++ ) {
if( pcr[pc2[j]] ) kn[m] |= bigbyte[j];
if( pcr[pc2[j+24]] ) kn[n] |= bigbyte[j];
}
}
cookey(kn);
return;
}

static void cookey(register unsigned long *raw1)
{
register unsigned long *cook, *raw0;
unsigned long dough[32];
register int i;

cook = dough;
for( i = 0; i < 16; i++, raw1++ ) {
raw0 = raw1++;
*cook	 = (*raw0 & 0x00fc0000L) << 6;
*cook	|= (*raw0 & 0x00000fc0L) << 10;
*cook	|= (*raw1 & 0x00fc0000L) >> 10;
*cook++ |= (*raw1 & 0x00000fc0L) >> 6;
*cook	 = (*raw0 & 0x0003f000L) << 12;
*cook	|= (*raw0 & 0x0000003fL) << 16;
*cook	|= (*raw1 & 0x0003f000L) >> 4;
*cook++ |= (*raw1 & 0x0000003fL);
}
usekey(dough);
return;
}

void cpkey(register unsigned long *into)
{
register unsigned long *from, *endp;

from = KnL, endp = &KnL[32];
while( from < endp ) *into++ = *from++;
return;
}

void usekey(register unsigned long *from)
{
register unsigned long *to, *endp;

to = KnL, endp = &KnL[32];
while( to < endp ) *to++ = *from++;
return;
}

void des(unsigned char *inblock,unsigned char *outblock)
{
unsigned long work[2];

scrunch(inblock, work);
desfunc(work, KnL);
unscrun(work, outblock);
return;
}

static void scrunch(register unsigned char *outof,register unsigned long *into)
{
*into	 = (*outof++ & 0xffL) << 24;
*into	|= (*outof++ & 0xffL) << 16;
*into	|= (*outof++ & 0xffL) << 8;
*into++ |= (*outof++ & 0xffL);
*into	 = (*outof++ & 0xffL) << 24;
*into	|= (*outof++ & 0xffL) << 16;
*into	|= (*outof++ & 0xffL) << 8;
*into	|= (*outof   & 0xffL);
return;
}

static void unscrun(register unsigned long *outof,register unsigned char *into)
{
*into++ = (*outof >> 24) & 0xffL;
*into++ = (*outof >> 16) & 0xffL;
*into++ = (*outof >>  8) & 0xffL;
*into++ =  *outof++	 & 0xffL;
*into++ = (*outof >> 24) & 0xffL;
*into++ = (*outof >> 16) & 0xffL;
*into++ = (*outof >>  8) & 0xffL;
*into	=  *outof	 & 0xffL;
return;
}

static unsigned long SP1[64] = {
0x01010400L, 0x00000000L, 0x00010000L, 0x01010404L,
0x01010004L, 0x00010404L, 0x00000004L, 0x00010000L,
0x00000400L, 0x01010400L, 0x01010404L, 0x00000400L,
0x01000404L, 0x01010004L, 0x01000000L, 0x00000004L,
0x00000404L, 0x01000400L, 0x01000400L, 0x00010400L,
0x00010400L, 0x01010000L, 0x01010000L, 0x01000404L,
0x00010004L, 0x01000004L, 0x01000004L, 0x00010004L,
0x00000000L, 0x00000404L, 0x00010404L, 0x01000000L,
0x00010000L, 0x01010404L, 0x00000004L, 0x01010000L,
0x01010400L, 0x01000000L, 0x01000000L, 0x00000400L,
0x01010004L, 0x00010000L, 0x00010400L, 0x01000004L,
0x00000400L, 0x00000004L, 0x01000404L, 0x00010404L,
0x01010404L, 0x00010004L, 0x01010000L, 0x01000404L,
0x01000004L, 0x00000404L, 0x00010404L, 0x01010400L,
0x00000404L, 0x01000400L, 0x01000400L, 0x00000000L,
0x00010004L, 0x00010400L, 0x00000000L, 0x01010004L };

static unsigned long SP2[64] = {
0x80108020L, 0x80008000L, 0x00008000L, 0x00108020L,
0x00100000L, 0x00000020L, 0x80100020L, 0x80008020L,
0x80000020L, 0x80108020L, 0x80108000L, 0x80000000L,
0x80008000L, 0x00100000L, 0x00000020L, 0x80100020L,
0x00108000L, 0x00100020L, 0x80008020L, 0x00000000L,
0x80000000L, 0x00008000L, 0x00108020L, 0x80100000L,
0x00100020L, 0x80000020L, 0x00000000L, 0x00108000L,
0x00008020L, 0x80108000L, 0x80100000L, 0x00008020L,
0x00000000L, 0x00108020L, 0x80100020L, 0x00100000L,
0x80008020L, 0x80100000L, 0x80108000L, 0x00008000L,
0x80100000L, 0x80008000L, 0x00000020L, 0x80108020L,
0x00108020L, 0x00000020L, 0x00008000L, 0x80000000L,
0x00008020L, 0x80108000L, 0x00100000L, 0x80000020L,
0x00100020L, 0x80008020L, 0x80000020L, 0x00100020L,
0x00108000L, 0x00000000L, 0x80008000L, 0x00008020L,
0x80000000L, 0x80100020L, 0x80108020L, 0x00108000L };

static unsigned long SP3[64] = {
0x00000208L, 0x08020200L, 0x00000000L, 0x08020008L,
0x08000200L, 0x00000000L, 0x00020208L, 0x08000200L,
0x00020008L, 0x08000008L, 0x08000008L, 0x00020000L,
0x08020208L, 0x00020008L, 0x08020000L, 0x00000208L,
0x08000000L, 0x00000008L, 0x08020200L, 0x00000200L,
0x00020200L, 0x08020000L, 0x08020008L, 0x00020208L,
0x08000208L, 0x00020200L, 0x00020000L, 0x08000208L,
0x00000008L, 0x08020208L, 0x00000200L, 0x08000000L,
0x08020200L, 0x08000000L, 0x00020008L, 0x00000208L,
0x00020000L, 0x08020200L, 0x08000200L, 0x00000000L,
0x00000200L, 0x00020008L, 0x08020208L, 0x08000200L,
0x08000008L, 0x00000200L, 0x00000000L, 0x08020008L,
0x08000208L, 0x00020000L, 0x08000000L, 0x08020208L,
0x00000008L, 0x00020208L, 0x00020200L, 0x08000008L,
0x08020000L, 0x08000208L, 0x00000208L, 0x08020000L,
0x00020208L, 0x00000008L, 0x08020008L, 0x00020200L };

static unsigned long SP4[64] = {
0x00802001L, 0x00002081L, 0x00002081L, 0x00000080L,
0x00802080L, 0x00800081L, 0x00800001L, 0x00002001L,
0x00000000L, 0x00802000L, 0x00802000L, 0x00802081L,
0x00000081L, 0x00000000L, 0x00800080L, 0x00800001L,
0x00000001L, 0x00002000L, 0x00800000L, 0x00802001L,
0x00000080L, 0x00800000L, 0x00002001L, 0x00002080L,
0x00800081L, 0x00000001L, 0x00002080L, 0x00800080L,
0x00002000L, 0x00802080L, 0x00802081L, 0x00000081L,
0x00800080L, 0x00800001L, 0x00802000L, 0x00802081L,
0x00000081L, 0x00000000L, 0x00000000L, 0x00802000L,
0x00002080L, 0x00800080L, 0x00800081L, 0x00000001L,
0x00802001L, 0x00002081L, 0x00002081L, 0x00000080L,
0x00802081L, 0x00000081L, 0x00000001L, 0x00002000L,
0x00800001L, 0x00002001L, 0x00802080L, 0x00800081L,
0x00002001L, 0x00002080L, 0x00800000L, 0x00802001L,
0x00000080L, 0x00800000L, 0x00002000L, 0x00802080L };

static unsigned long SP5[64] = {
0x00000100L, 0x02080100L, 0x02080000L, 0x42000100L,
0x00080000L, 0x00000100L, 0x40000000L, 0x02080000L,
0x40080100L, 0x00080000L, 0x02000100L, 0x40080100L,
0x42000100L, 0x42080000L, 0x00080100L, 0x40000000L,
0x02000000L, 0x40080000L, 0x40080000L, 0x00000000L,
0x40000100L, 0x42080100L, 0x42080100L, 0x02000100L,
0x42080000L, 0x40000100L, 0x00000000L, 0x42000000L,
0x02080100L, 0x02000000L, 0x42000000L, 0x00080100L,
0x00080000L, 0x42000100L, 0x00000100L, 0x02000000L,
0x40000000L, 0x02080000L, 0x42000100L, 0x40080100L,
0x02000100L, 0x40000000L, 0x42080000L, 0x02080100L,
0x40080100L, 0x00000100L, 0x02000000L, 0x42080000L,
0x42080100L, 0x00080100L, 0x42000000L, 0x42080100L,
0x02080000L, 0x00000000L, 0x40080000L, 0x42000000L,
0x00080100L, 0x02000100L, 0x40000100L, 0x00080000L,
0x00000000L, 0x40080000L, 0x02080100L, 0x40000100L };

static unsigned long SP6[64] = {
0x20000010L, 0x20400000L, 0x00004000L, 0x20404010L,
0x20400000L, 0x00000010L, 0x20404010L, 0x00400000L,
0x20004000L, 0x00404010L, 0x00400000L, 0x20000010L,
0x00400010L, 0x20004000L, 0x20000000L, 0x00004010L,
0x00000000L, 0x00400010L, 0x20004010L, 0x00004000L,
0x00404000L, 0x20004010L, 0x00000010L, 0x20400010L,
0x20400010L, 0x00000000L, 0x00404010L, 0x20404000L,
0x00004010L, 0x00404000L, 0x20404000L, 0x20000000L,
0x20004000L, 0x00000010L, 0x20400010L, 0x00404000L,
0x20404010L, 0x00400000L, 0x00004010L, 0x20000010L,
0x00400000L, 0x20004000L, 0x20000000L, 0x00004010L,
0x20000010L, 0x20404010L, 0x00404000L, 0x20400000L,
0x00404010L, 0x20404000L, 0x00000000L, 0x20400010L,
0x00000010L, 0x00004000L, 0x20400000L, 0x00404010L,
0x00004000L, 0x00400010L, 0x20004010L, 0x00000000L,
0x20404000L, 0x20000000L, 0x00400010L, 0x20004010L };

static unsigned long SP7[64] = {
0x00200000L, 0x04200002L, 0x04000802L, 0x00000000L,
0x00000800L, 0x04000802L, 0x00200802L, 0x04200800L,
0x04200802L, 0x00200000L, 0x00000000L, 0x04000002L,
0x00000002L, 0x04000000L, 0x04200002L, 0x00000802L,
0x04000800L, 0x00200802L, 0x00200002L, 0x04000800L,
0x04000002L, 0x04200000L, 0x04200800L, 0x00200002L,
0x04200000L, 0x00000800L, 0x00000802L, 0x04200802L,
0x00200800L, 0x00000002L, 0x04000000L, 0x00200800L,
0x04000000L, 0x00200800L, 0x00200000L, 0x04000802L,
0x04000802L, 0x04200002L, 0x04200002L, 0x00000002L,
0x00200002L, 0x04000000L, 0x04000800L, 0x00200000L,
0x04200800L, 0x00000802L, 0x00200802L, 0x04200800L,
0x00000802L, 0x04000002L, 0x04200802L, 0x04200000L,
0x00200800L, 0x00000000L, 0x00000002L, 0x04200802L,
0x00000000L, 0x00200802L, 0x04200000L, 0x00000800L,
0x04000002L, 0x04000800L, 0x00000800L, 0x00200002L };

static unsigned long SP8[64] = {
0x10001040L, 0x00001000L, 0x00040000L, 0x10041040L,
0x10000000L, 0x10001040L, 0x00000040L, 0x10000000L,
0x00040040L, 0x10040000L, 0x10041040L, 0x00041000L,
0x10041000L, 0x00041040L, 0x00001000L, 0x00000040L,
0x10040000L, 0x10000040L, 0x10001000L, 0x00001040L,
0x00041000L, 0x00040040L, 0x10040040L, 0x10041000L,
0x00001040L, 0x00000000L, 0x00000000L, 0x10040040L,
0x10000040L, 0x10001000L, 0x00041040L, 0x00040000L,
0x00041040L, 0x00040000L, 0x10041000L, 0x00001000L,
0x00000040L, 0x10040040L, 0x00001000L, 0x00041040L,
0x10001000L, 0x00000040L, 0x10000040L, 0x10040000L,
0x10040040L, 0x10000000L, 0x00040000L, 0x10001040L,
0x00000000L, 0x10041040L, 0x00040040L, 0x10000040L,
0x10040000L, 0x10001000L, 0x10001040L, 0x00000000L,
0x10041040L, 0x00041000L, 0x00041000L, 0x00001040L,
0x00001040L, 0x00040040L, 0x10000000L, 0x10041000L };

static void desfunc(register unsigned long *block,register unsigned long *keys)
{
register unsigned long fval, work, right, leftt;
register int round;

leftt = block[0];
right = block[1];
work = ((leftt >> 4) ^ right) & 0x0f0f0f0fL;
right ^= work;
leftt ^= (work << 4);
work = ((leftt >> 16) ^ right) & 0x0000ffffL;
right ^= work;
leftt ^= (work << 16);
work = ((right >> 2) ^ leftt) & 0x33333333L;
leftt ^= work;
right ^= (work << 2);
work = ((right >> 8) ^ leftt) & 0x00ff00ffL;
leftt ^= work;
right ^= (work << 8);
right = ((right << 1) | ((right >> 31) & 1L)) & 0xffffffffL;
work = (leftt ^ right) & 0xaaaaaaaaL;
leftt ^= work;
right ^= work;
leftt = ((leftt << 1) | ((leftt >> 31) & 1L)) & 0xffffffffL;

for( round = 0; round < 8; round++ ) {
work  = (right << 28) | (right >> 4);
work ^= *keys++;
fval  = SP7[ work		 & 0x3fL];
fval |= SP5[(work >>  8) & 0x3fL];
fval |= SP3[(work >> 16) & 0x3fL];
fval |= SP1[(work >> 24) & 0x3fL];
work  = right ^ *keys++;
fval |= SP8[ work		 & 0x3fL];
fval |= SP6[(work >>  8) & 0x3fL];
fval |= SP4[(work >> 16) & 0x3fL];
fval |= SP2[(work >> 24) & 0x3fL];
leftt ^= fval;
work  = (leftt << 28) | (leftt >> 4);
work ^= *keys++;
fval  = SP7[ work		 & 0x3fL];
fval |= SP5[(work >>  8) & 0x3fL];
fval |= SP3[(work >> 16) & 0x3fL];
fval |= SP1[(work >> 24) & 0x3fL];
work  = leftt ^ *keys++;
fval |= SP8[ work		 & 0x3fL];
fval |= SP6[(work >>  8) & 0x3fL];
fval |= SP4[(work >> 16) & 0x3fL];
fval |= SP2[(work >> 24) & 0x3fL];
right ^= fval;
}

right = (right << 31) | (right >> 1);
work = (leftt ^ right) & 0xaaaaaaaaL;
leftt ^= work;
right ^= work;
leftt = (leftt << 31) | (leftt >> 1);
work = ((leftt >> 8) ^ right) & 0x00ff00ffL;
right ^= work;
leftt ^= (work << 8);
work = ((leftt >> 2) ^ right) & 0x33333333L;
right ^= work;
leftt ^= (work << 2);
work = ((right >> 16) ^ leftt) & 0x0000ffffL;
leftt ^= work;
right ^= (work << 16);
work = ((right >> 4) ^ leftt) & 0x0f0f0f0fL;
leftt ^= work;
right ^= (work << 4);
*block++ = right;
*block = leftt;
return;
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





int enc(const char* pass, unsigned char *inbuf, int inlen, unsigned char *outbuf)
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
	int i = 0;
	for (i=0; i<step; i++)
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
	return 0;
}

int dec(const char* pass, unsigned char *inbuf, int inlen, unsigned char *outbuf, int outlen)
{
	unsigned char passwd[8] = {0};
	StrToKey(pass, (char*)passwd);
	//	if (strlen(pass) > 8)
	//		memcpy(passwd, pass, 8);
	//	else
	//		memcpy(passwd, pass, strlen(pass));


	if (inlen % 8 != 2)
	{
		return -1;
	}

	WORD size;
	memcpy(&size, inbuf, sizeof(WORD));
	ASSERT(size + 2 <= inlen);
	outlen = size;

	int step = (inlen - 2) / 8;

	deskey(passwd, DE1);
	int i = 0;
	for (i=0; i<step; i++)
	{
		des(&inbuf[2+8*i], &outbuf[8*i]);
	}
	return outlen;
}


int enc_len(int len)
{
	if (len % 8 == 0)
		return len + 2;
	else
		return (len / 8 + 1)*8 + 2;
}

*/
import "C"
import (
	"fmt"
	"os"
	"bufio"
	"io"
)
/*
func EncLen(len int) int{
	return int(C.enc_len(C.int(len)))
}
*/
func EncLen(len int) int{
	if len % 8 == 0 {
		return len + 2
	}else{
		return (len / 8 + 1) * 8 +2
	}
}

//int enc(const char* pass, unsigned char *inbuf, int inlen, unsigned char *outbuf)
func Encrypt(inbuf []byte,outbuf []byte)(outbyte []byte,err error){
	passptr := C.CString(DES_KEY)
	inbufptr := C.CBytes(inbuf)
	outbufptr  := C.CBytes(outbuf)
	inlen := len(inbuf)
	ret := C.enc(passptr,(*C.uchar)(inbufptr),C.int(inlen),(*C.uchar)(outbufptr))

	if ret == 0{
		outlen := EncLen(inlen)
		outbuf = C.GoBytes(outbufptr,C.int(outlen))
		return outbuf[:outlen],nil
	}else {
		return nil,fmt.Errorf("encrypt fail")
	}
}

//int dec(const char* pass, unsigned char *inbuf, int inlen, unsigned char *outbuf, int outlen)
func Decrypt(inbuf []byte,outbuf []byte)( []byte, error){
	passptr := C.CString(DES_KEY)
	inbufptr := C.CBytes(inbuf)
	outbufptr  := C.CBytes(outbuf)
	inlen := len(inbuf)
	var outlen C.int
	outlen = C.dec(passptr,(*C.uchar)(inbufptr),C.int(inlen),(*C.uchar)(outbufptr),outlen)

	length := int(outlen)
	if length != -1{
		outbuf = C.GoBytes(outbufptr,outlen)
		return outbuf[:length],nil
	}else {
		return nil,fmt.Errorf("decrypt fail")
	}

}


func EncFile(srcFile, dstFile string)error  {
	srcfile, errSrc := os.Open(srcFile)
	if errSrc != nil { return errSrc}
	srcfileRead := bufio.NewReader(srcfile)


	dstfile, errDst := os.Create(dstFile)
	if errDst != nil {return  errDst}
	dstfileWrite := bufio.NewWriter(dstfile)

	defer srcfile.Close()
	defer dstfile.Close()

	srcBuf := make([]byte, MAX_DATA_LEN)
	for {
		n, errRead := srcfileRead.Read(srcBuf)
		if errRead != nil && errRead != io.EOF {return errRead}
		if 0 == n {break}
		data, errMake := MakeDataPacket(srcBuf[:n])
		if errMake != nil {return errMake}
		wn,errWrite := dstfileWrite.Write(data)
		if errWrite != nil {return errWrite}
		if wn != len(data) {return fmt.Errorf("write data len is wrong %v",wn)}
		if err := dstfileWrite.Flush(); err != nil {return err}
	}
	return nil
}
