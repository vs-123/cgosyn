#include "stdio.h"
#include "math.h"

int main()
{
	int age = 21;
	if (isAdult(age))
	{	
		printf("You are an adult!");
	}
	else
	{
		printf("You are younger than 18!");
	}
	return 0;
}

int isAdult(int age)
{
	return age >= 18;
}
