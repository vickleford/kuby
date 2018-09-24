#include <curl/curl.h>

class CurlTest {
    const char *url;
public:
    CurlTest(const char *url);
    void get();
};
