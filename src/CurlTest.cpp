#include <iostream>
#include "CurlTest.h"

CurlTest::CurlTest(const char *url) : url{url} {}

void CurlTest::get() {
    std::cout << "got " << url << std::endl;
}
