#include "vec.h"


// data_t dotproduct(vec_ptr u, vec_ptr v) {
//    data_t sum = 0, u_val, v_val;
//
//    for (long i = 0; i < vec_length(u); i++) { // we can assume both vectors are same length
//         get_vec_element(u, i, &u_val);
//         get_vec_element(v, i, &v_val);
//         sum += u_val * v_val;
//    }
//    return sum;
// }


// data_t dotproduct(vec_ptr u, vec_ptr v) {
//    data_t sum = 0, *u_data = get_vec_start(u), *v_data = get_vec_start(v);
//    int j = vec_length(u);
//    long i = 0;
//    for (; i < j; i+=2) {
//         long p1 = u_data[i] * v_data[i];
//         long p2 = u_data[i+1] * v_data[i+1];
//         sum += p1 + p2;
//    }
//    for (; i < j; i++) {
//         sum += u_data[i] * v_data[i];
//    }
//    return sum;
// }

// data_t dotproduct(vec_ptr u, vec_ptr v) {
//   data_t *u_data = get_vec_start(u), *v_data = get_vec_start(v), sum1 = 0, sum2 = 0, sum3 = 0, sum4 = 0;
//   long i, n = vec_length(u);
//
//   for (i = 0; i < n - 3; i += 4) {
//     sum1 += u_data[i] * v_data[i];
//     sum2 += u_data[i + 1] * v_data[i + 1];
//     sum3 += u_data[i + 2] * v_data[i + 2];
//     sum4 += u_data[i + 3] * v_data[i + 3];
//   }
//
//   for (; i < n; i++) {
//     sum1 += u_data[i] * v_data[i];
//   }
//   return sum1 + sum2 + sum3 + sum4;
// }

data_t dotproduct(vec_ptr u, vec_ptr v) {
  data_t sum1 = 0, sum2 = 0, sum3 = 0, sum4 = 0, *u_data = get_vec_start(u),
         *v_data = get_vec_start(v);
  long i, n = vec_length(u);

  for (i = 0; i < n - 3; i += 4) {
    sum1 += u_data[i] * v_data[i];
    sum2 += u_data[i + 1] * v_data[i + 1];
    sum3 += u_data[i + 2] * v_data[i + 2];
    sum4 += u_data[i + 3] * v_data[i + 3];
  }

  for (; i < n; i++) {
    sum1 += u_data[i] * v_data[i];
  }
  return sum1 + sum2 + sum3 + sum4;
}
