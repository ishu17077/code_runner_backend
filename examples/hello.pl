#!/usr/bin/perl
use strict;
use warnings;
my $input_line = <STDIN>;
chomp($input_line);
my $num = int($input_line);

if ($num % 2 == 0){
    print("Yes\n");
}else{
    print("No\n");
}