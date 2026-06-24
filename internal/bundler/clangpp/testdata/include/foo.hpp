#ifndef FOO1_HPP
#define FOO1_HPP

#include <iostream>

#include "bar.hpp"
#include "baz.hpp"

void foo1() {
  bar();
  baz();
  std::cout << "foo1" << std::endl;
}

#endif  // FOO1_HPP
