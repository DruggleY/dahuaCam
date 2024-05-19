#ifndef DHNETSDK_H
#define DHNETSDK_H



#define LLONG       long
#define INT64       long long
#define TP_U64      unsigned long long
#define LDWORD      long
#define WORD        unsigned short
#define DWORD       unsigned int
#define LONG        int
#define BYTE        unsigned char
#define UINT        unsigned int
#define HDC         void*
#define HWND        void*
#define LPVOID      void*
#define BOOL    int
#define CALLBACK


#define DH_SERIALNO_LEN                   48


typedef void (CALLBACK *fDisConnect)(LLONG lLoginID, char *pchDVRIP, LONG nDVRPort, LDWORD dwUser);
typedef void (CALLBACK *fSnapRev)(LLONG lLoginID, BYTE *pBuf, UINT RevLen, UINT EncodeType, DWORD CmdSerial, LDWORD dwUser);

typedef enum tagEM_LOGIN_SPAC_CAP_TYPE
{
    EM_LOGIN_SPEC_CAP_TCP               = 0,    /// TCP��½, Ĭ�Ϸ�ʽ
    EM_LOGIN_SPEC_CAP_ANY               = 1,    /// ��������½
    EM_LOGIN_SPEC_CAP_SERVER_CONN       = 2,    /// ����ע��ĵ���
    EM_LOGIN_SPEC_CAP_MULTICAST         = 3,    /// �鲥��½
    EM_LOGIN_SPEC_CAP_UDP               = 4,    /// UDP��ʽ�µĵ���
    EM_LOGIN_SPEC_CAP_MAIN_CONN_ONLY    = 6,    /// ֻ���������µĵ���
    EM_LOGIN_SPEC_CAP_SSL               = 7,    /// SSL���ܷ�ʽ��½

    EM_LOGIN_SPEC_CAP_INTELLIGENT_BOX   = 9,    /// ��¼���ܺ�Զ���豸
    EM_LOGIN_SPEC_CAP_NO_CONFIG         = 10,   /// ��¼�豸����ȡ���ò���
    EM_LOGIN_SPEC_CAP_U_LOGIN           = 11,   /// ��U���豸�ĵ���
    EM_LOGIN_SPEC_CAP_LDAP              = 12,   /// LDAP��ʽ��¼
    EM_LOGIN_SPEC_CAP_AD                = 13,   /// AD��ActiveDirectory����¼��ʽ
    EM_LOGIN_SPEC_CAP_RADIUS            = 14,   /// Radius ��¼��ʽ
    EM_LOGIN_SPEC_CAP_SOCKET_5          = 15,   /// Socks5��½��ʽ
    EM_LOGIN_SPEC_CAP_CLOUD             = 16,   /// �Ƶ�½��ʽ
    EM_LOGIN_SPEC_CAP_AUTH_TWICE        = 17,   /// ���μ�Ȩ��½��ʽ
    EM_LOGIN_SPEC_CAP_TS                = 18,   /// TS�����ͻ��˵�½��ʽ
    EM_LOGIN_SPEC_CAP_P2P               = 19,   /// ΪP2P��½��ʽ
    EM_LOGIN_SPEC_CAP_MOBILE            = 20,   /// �ֻ��ͻ��˵�½
	EM_LOGIN_SPEC_CAP_TLS_ADAPTER       = 21,   /// ����Ӧtls����
	EM_LOGIN_SPEC_CAP_TLS_COMPEL        = 22,   /// ǿ��tls����
	EM_LOGIN_SPEC_CAP_TLS_MAIN_ONLY     = 23,   /// ����tls����
	EM_LOGIN_SPEC_CAP_NTLS_VERIFY       = 24,   /// ����tls˫����֤��¼-��׼��ȫ��·TLS,�����ȵ���CLIENT_SetSDKLocalCfg�ӿ�,���ù���TLS��ص�֤���˽Կ�ļ�·��
    EM_LOGIN_SPEC_CAP_INVALID                   /// ��Ч�ĵ�½��ʽ
}EM_LOGIN_SPAC_CAP_TYPE;

///@brief ��½ʱTLS����ģʽ
typedef enum tagEM_LOGIN_TLS_TYPE
{
	EM_LOGIN_TLS_TYPE_NO_TLS            = 0,    /// ����ģʽ- ����tls����, Ĭ�Ϸ�ʽ
	EM_LOGIN_TLS_TYPE_TLS_ADAPTER       = 1,   /// ��֯ģʽ-����Ӧtls����
	EM_LOGIN_TLS_TYPE_TLS_COMPEL        = 2,   /// ��֯ģʽ-ǿ��tls����
	EM_LOGIN_TLS_TYPE_TLS_MAIN_ONLY     = 3,   /// ��֯ģʽ-����tls����
	EM_LOGIN_TLS_TYPE_TLS_GENERAL       = 4,   /// ��׼ģʽ-ͨ��tls���ܣ�������·ͬһ�˿�
	EM_LOGIN_TLS_TYPE_TLS_UPNP       = 5,		  /// ��׼ģʽ-ͨ��tls���ܣ�������·��ͬ�˿ڣ�ͬʱ�˿�ӳ���֧��UPNPģʽ��tls���ܣ�������ǿ�Ƶ�ͨ��TLS,�����Ӳ���TLS
}EM_LOGIN_TLS_TYPE;

///@brief ץͼ�����ṹ��
typedef struct _snap_param
{
    unsigned int     Channel;                       /// ץͼ��ͨ��
    unsigned int     Quality;                       /// ���ʣ�1~6
    unsigned int     ImageSize;                     /// �����С��0��QCIF,1��CIF,2��D1
    unsigned int     mode;                          /// ץͼģʽ��-1:��ʾֹͣץͼ, 0����ʾ����һ֡, 1����ʾ��ʱ��������, 2����ʾ��������
    unsigned int     InterSnap;                     /// ʱ�䵥λ�룻��mode=1��ʾ��ʱ��������ʱ
													/// ֻ�в��������豸(�磺�����豸)֧��ͨ�����ֶ�ʵ�ֶ�ʱץͼʱ����������
													/// ����ͨ�� CFG_CMD_ENCODE ���õ�stuSnapFormat[nSnapMode].stuVideoFormat.nFrameRate�ֶ�ʵ����ع���
    unsigned int     CmdSerial;                     /// �������кţ���Чֵ��Χ 0~65535��������Χ�ᱻ�ض�Ϊ unsigned short
    unsigned int     Reserved[4];
} SNAP_PARAMS, *LPSNAP_PARAMS;

///@brief �豸��Ϣ��չ
typedef struct
{
    BYTE                sSerialNumber[DH_SERIALNO_LEN];     /// ���к�
    int                 nAlarmInPortNum;                    /// DVR�����������
    int                 nAlarmOutPortNum;                   /// DVR�����������
    int                 nDiskNum;                           /// DVRӲ�̸���
    int                 nDVRType;                           /// DVR����,��ö�� NET_DEVICE_TYPE
    int                 nChanNum;                           /// DVRͨ������
    BYTE                byLimitLoginTime;                   /// ���߳�ʱʱ��,Ϊ0��ʾ�����Ƶ�½,��0��ʾ���Ƶķ�����
    BYTE                byLeftLogTimes;                     /// ����½ʧ��ԭ��Ϊ�������ʱ,ͨ���˲���֪ͨ�û�,ʣ���½����,Ϊ0ʱ��ʾ�˲�����Ч
    BYTE                bReserved[2];                       /// �����ֽ�,,�ڲ���ʹ��
    int                 nLockLeftTime;                      /// ����½ʧ��,�û�����ʣ��ʱ�䣨������, -1��ʾ�豸δ���øò���
    char                Reserved[4];                        /// ����,�ڲ���ʹ��
    int                 nNTlsPort;                          /// ����TLS��¼�˿�,����¼������Ϊ24ʱ��Ч
    char                Reserved2[16];                      /// ����
} NET_DEVICEINFO_Ex, *LPNET_DEVICEINFO_Ex;

typedef struct tagNET_IN_LOGIN_WITH_HIGHLEVEL_SECURITY
{
	DWORD						dwSize;				/// �ṹ���С
	char						szIP[64];			/// IP
	int							nPort;				/// �˿�
	char						szUserName[64];		/// �û���
	char						szPassword[64];		/// ����
	EM_LOGIN_SPAC_CAP_TYPE		emSpecCap;			/// ��¼ģʽ
	BYTE						byReserved[4];		/// �ֽڶ���
	void*						pCapParam;			/// �� CLIENT_LoginEx �ӿ� pCapParam �� nSpecCap ��ϵ
	EM_LOGIN_TLS_TYPE           emTLSCap;           /// ��¼��TLSģʽ��Ŀǰ��֧��emSpecCapΪEM_LOGIN_SPEC_CAP_TCP��EM_LOGIN_SPEC_CAP_SERVER_CONN ģʽ�µ� tls��½(TLS��������ʹ�ø�ѡ��)
}NET_IN_LOGIN_WITH_HIGHLEVEL_SECURITY;

///@brief CLIENT_LoginWithHighLevelSecurity �������
typedef struct tagNET_OUT_LOGIN_WITH_HIGHLEVEL_SECURITY
{
	DWORD						dwSize;				/// �ṹ���С
	NET_DEVICEINFO_Ex			stuDeviceInfo;		/// �豸��Ϣ
	int							nError;				/// �����룬�� CLIENT_Login �ӿڴ�����
	BYTE						byReserved[132];	/// Ԥ���ֶ�
}NET_OUT_LOGIN_WITH_HIGHLEVEL_SECURITY;

//BOOL CLIENT_Init(fDisConnect cbDisConnect, LDWORD dwUser);
//LLONG CLIENT_LoginWithHighLevelSecurity(NET_IN_LOGIN_WITH_HIGHLEVEL_SECURITY* pstInParam, NET_OUT_LOGIN_WITH_HIGHLEVEL_SECURITY* pstOutParam);
//void CLIENT_SetSnapRevCallBack(fSnapRev OnSnapRevMessage, LDWORD dwUser);
//BOOL CLIENT_SnapPictureEx(LLONG lLoginID, SNAP_PARAMS *par, int *reserved);
//BOOL CLIENT_Logout(LLONG lLoginID);

#endif