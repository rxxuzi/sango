#define _GNU_SOURCE  // for strdup on Linux
#include "sango.h"
#include <stdarg.h>
#include <stddef.h>
#include <stdint.h>

// String duplication helper for portability
static char* sango_strdup(const char* s) {
    size_t len = strlen(s) + 1;
    char* dup = (char*)malloc(len);
    if (dup) {
        memcpy(dup, s, len);
    }
    return dup;
}

// String operations
sango_string sango_string_concat(sango_string s1, sango_string s2) {
    size_t len1 = strlen(s1);
    size_t len2 = strlen(s2);
    sango_string result = (sango_string)malloc(len1 + len2 + 1);
    strcpy(result, s1);
    strcat(result, s2);
    return result;
}

sango_string sango_string_repeat(sango_string s, sango_int count) {
    size_t len = strlen(s);
    size_t total_len = len * count;
    sango_string result = (sango_string)malloc(total_len + 1);
    result[0] = '\0';
    for (int i = 0; i < count; i++) {
        strcat(result, s);
    }
    return result;
}

sango_string sango_string_from_int(sango_int n) {
    char buffer[32];
    snprintf(buffer, sizeof(buffer), "%d", n);
    return sango_strdup(buffer);
}

sango_string sango_string_from_long(sango_long n) {
    char buffer[32];
    snprintf(buffer, sizeof(buffer), "%lld", (long long)n);
    return sango_strdup(buffer);
}

sango_string sango_string_from_float(sango_float f) {
    char buffer[64];
    snprintf(buffer, sizeof(buffer), "%f", f);
    return sango_strdup(buffer);
}

sango_string sango_string_from_double(sango_double d) {
    char buffer[64];
    snprintf(buffer, sizeof(buffer), "%lf", d);
    return sango_strdup(buffer);
}

sango_int sango_string_to_int(sango_string s) {
    return (sango_int)atoi(s);
}

sango_long sango_string_to_long(sango_string s) {
    return (sango_long)atoll(s);
}

sango_float sango_string_to_float(sango_string s) {
    return (sango_float)atof(s);
}

sango_double sango_string_to_double(sango_string s) {
    return (sango_double)atof(s);
}

// Built-in functions
void sango_print(const char* format, ...) {
    va_list args;
    va_start(args, format);
    vprintf(format, args);
    va_end(args);
}

void sango_println(const char* format, ...) {
    va_list args;
    va_start(args, format);
    vprintf(format, args);
    va_end(args);
    printf("\n");
}

size_t sango_len_string(sango_string s) {
    return strlen(s);
}

size_t sango_len_array(void* arr) {
    sango_array* array = (sango_array*)arr;
    return array->length;
}

void sango_assert(sango_bool condition, const char* message) {
    if (!condition) {
        fprintf(stderr, "Assertion failed: %s\n", message);
        abort();
    }
}

void sango_panic(const char* message) {
    fprintf(stderr, "Panic: %s\n", message);
    abort();
}

// Memory management
void* sango_alloc(size_t size) {
    void* ptr = malloc(size);
    if (!ptr) {
        sango_panic("Out of memory");
    }
    return ptr;
}

void sango_free(void* ptr) {
    free(ptr);
}

// Array implementation
sango_array* sango_array_new(size_t element_size, size_t initial_capacity) {
    sango_array* arr = (sango_array*)sango_alloc(sizeof(sango_array));
    arr->element_size = element_size;
    arr->length = 0;
    arr->capacity = initial_capacity > 0 ? initial_capacity : 16;
    arr->data = sango_alloc(arr->capacity * element_size);
    return arr;
}

void sango_array_free(sango_array* arr) {
    if (arr) {
        sango_free(arr->data);
        sango_free(arr);
    }
}

void sango_array_push(sango_array* arr, void* element) {
    if (arr->length >= arr->capacity) {
        arr->capacity *= 2;
        arr->data = realloc(arr->data, arr->capacity * arr->element_size);
        if (!arr->data) {
            sango_panic("Out of memory");
        }
    }
    memcpy((char*)arr->data + arr->length * arr->element_size, element, arr->element_size);
    arr->length++;
}

void* sango_array_get(sango_array* arr, size_t index) {
    if (index >= arr->length) {
        sango_panic("Array index out of bounds");
    }
    return (char*)arr->data + index * arr->element_size;
}

sango_array* sango_array_slice(sango_array* arr, size_t start, size_t end) {
    if (start > end || end > arr->length) {
        sango_panic("Invalid slice range");
    }
    size_t new_length = end - start;
    sango_array* new_arr = sango_array_new(arr->element_size, new_length);
    new_arr->length = new_length;
    memcpy(new_arr->data, (char*)arr->data + start * arr->element_size, 
           new_length * arr->element_size);
    return new_arr;
}

sango_array* sango_array_concat(sango_array* arr1, sango_array* arr2) {
    if (arr1->element_size != arr2->element_size) {
        sango_panic("Cannot concat arrays of different types");
    }
    size_t total_length = arr1->length + arr2->length;
    sango_array* new_arr = sango_array_new(arr1->element_size, total_length);
    new_arr->length = total_length;
    memcpy(new_arr->data, arr1->data, arr1->length * arr1->element_size);
    memcpy((char*)new_arr->data + arr1->length * arr1->element_size, 
           arr2->data, arr2->length * arr2->element_size);
    return new_arr;
}