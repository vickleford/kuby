#include "CurlTest.h"

int main() {
    CurlTest c {"http://example.com"};
    c.get();
    return 0;
}
