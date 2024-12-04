#include <stdio.h>
#include <stdlib.h>

#define MUL_OPEN "mul("
#define MUL_CLOSE ')'
#define MUL_DIGIT_SEP ','

int main(int argc, char **argv)
{
    printf("Reading input file...\n");

    FILE *f = fopen("input.txt", "r");
    if (f == NULL) {
        fprintf(stderr, "Cannot read input\n");
        return 1;
    }

    printf("Parsing input file...\n");

    size_t sum = 0;
    int a = 0, b = 0;
    int mul_open_index = 0;

    for (char c; c = fgetc(f), c != EOF;)
    {
        if (mul_open_index == 0) {
            a = b = 0;
        }

        switch (mul_open_index) {
        // Reading MUL_OPEN
        case 0:
        case 1:
        case 2:
        case 3:
            if (c != MUL_OPEN[mul_open_index])
            {
                mul_open_index = 0;
            } else {
                mul_open_index++;
            }
            break;
        // Reading first int to MUL_DIGIT_SEP
        case 4:
            if (c == MUL_DIGIT_SEP) {
                mul_open_index++;
            } else if (c >= '0' && c <= '9') {
                a *= 10;
                a += c - '0';
            } else {
                mul_open_index = 0;
            }
            break;
        // Reading second int to MUL_CLOSE
        case 5:
            if (c == MUL_CLOSE) {
                sum += (a * b);
                mul_open_index = 0;
                printf("Read a=\t%d, b=\t%d, new sum is \t%lu\n", a, b, sum);
            } else if (c >= '0' && c <= '9') {
                b *= 10;
                b += c - '0';
            } else {
                mul_open_index = 0;
            }
            break;
        }
    }

    printf("Complete! output=%lu\n", sum);
    fclose(f);
}
