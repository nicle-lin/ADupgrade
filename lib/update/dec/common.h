#if !defined(__common_h)
#define __common_h

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
typedef void*		    PVOID;
typedef unsigned char   BOOLEAN;
typedef unsigned char   BOOL;
typedef void		    VOID;
typedef	unsigned int*	PULONG;
typedef long		    LONG ;
typedef char		    CHAR ;
typedef unsigned char	BYTE;

#define ASSERT assert


#define TRUE 1
#define FALSE 0
#endif
