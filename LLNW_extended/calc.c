float calc(double b2, float c1, double d2, float e1, double f2){
    #define K 0x00025630
    return (long double)K + (long double)b2 - (long double)d2/(long double)c1 + (long double)e1*(long double)f2;
}