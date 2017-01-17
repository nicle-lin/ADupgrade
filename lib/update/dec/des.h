/*
#ifdefined(__cplusplus)||defined(c_plusplus) //跨平台定义方法
extern "C"{
#endif
//... 正常的声明段
int enc_len(int len);
bool enc(const char* pass, BYTE* inbuf, int inlen, BYTE* outbuf);
bool dec(const char* pass, BYTE* inbuf, int inlen, BYTE* outbuf, int& outlen);
bool enc_char(const BYTE pass[8], BYTE* inbuf, int inlen, BYTE* outbuf);
bool dec_char(const BYTE pass[8], BYTE* inbuf, int inlen, BYTE* outbuf, int& outlen);
bool enc_passwd(const char* pass, BYTE data[16]);
bool dec_passwd(const char* pass, BYTE data[16]);
#ifdefined(__cplusplus)||defined(c_plusplus)
}
#endif
*/
//... 正常的声明段
extern int enc_len(int len);
extern bool enc(const char* pass, unsigned char* inbuf, int inlen, unsigned char* outbuf);
extern bool dec(const char* pass, unsigned char* inbuf, int inlen, unsigned char* outbuf, int &outlen);
extern bool enc_char(const unsigned char pass[8], unsigned char* inbuf, int inlen, unsigned char* outbuf);
extern bool dec_char(const unsigned char pass[8], unsigned char* inbuf, int inlen, unsigned char* outbuf, int &outlen);
extern bool enc_passwd(const char* pass, unsigned char data[16]);
extern bool dec_passwd(const char* pass, unsigned char data[16]);
