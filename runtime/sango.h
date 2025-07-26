#ifndef SANGO_H
#define SANGO_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>
#include <stdint.h>
#include <stddef.h>

// Sango basic types
typedef int32_t sango_int;
typedef int64_t sango_long;
typedef float sango_float;
typedef double sango_double;
typedef bool sango_bool;
typedef char* sango_string;

// Array structure definition (moved before function declarations)
typedef struct {
    void* data;
    size_t length;
    size_t capacity;
    size_t element_size;
} sango_array;

// String operations
sango_string sango_string_concat(sango_string s1, sango_string s2);
sango_string sango_string_repeat(sango_string s, sango_int count);
sango_string sango_string_from_int(sango_int n);
sango_string sango_string_from_long(sango_long n);
sango_string sango_string_from_float(sango_float f);
sango_string sango_string_from_double(sango_double d);
sango_int sango_string_to_int(sango_string s);
sango_long sango_string_to_long(sango_string s);
sango_float sango_string_to_float(sango_string s);
sango_double sango_string_to_double(sango_string s);

// Built-in functions
void sango_print(const char* format, ...);
void sango_println(const char* format, ...);
size_t sango_len_string(sango_string s);
size_t sango_len_array(void* arr);
void sango_assert(sango_bool condition, const char* message);
void sango_panic(const char* message) __attribute__((noreturn));

// Memory management helpers
void* sango_alloc(size_t size);
void sango_free(void* ptr);

// Array helpers
sango_array* sango_array_new(size_t element_size, size_t initial_capacity);
void sango_array_free(sango_array* arr);
void sango_array_push(sango_array* arr, void* element);
void* sango_array_get(sango_array* arr, size_t index);
sango_array* sango_array_slice(sango_array* arr, size_t start, size_t end);
sango_array* sango_array_concat(sango_array* arr1, sango_array* arr2);

#endif // SANGO_H