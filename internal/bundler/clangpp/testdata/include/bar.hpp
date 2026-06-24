#ifndef BAR_HPP
#define BAR_HPP

#include <iostream>

#include "qux.hpp"

void bar() {
  qux();
  std::cout << "bar" << std::endl;
}

#endif  // BAR_HPP
