 #include <stdio.h>
 #include <stdint.h>

 int main( int argc, char **argv ){
         uint16_t word = 0xff00;
         char *byte1 = (char *) &word;
         if( *byte1 ) printf( "This system is big-endian.\n" );
         else printf( "This system is little-endian.\n" );
         return 0;
}
