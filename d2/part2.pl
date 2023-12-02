#!/usr/bin/perl -w

use strict;

my $sum = 0;

sub max {
    my $m = shift;
    for $a (@_) {
        if ($a > $m) {
            $m = $a;
        }
    }
    return $m;
}

sub product {
    my $product = 1;
    for my $factor (@_) {
        $product *= $factor;
    }
    return $product;
}

foreach my $line (<ARGV>) {
    $line =~ /Game (\d+): (.*)/ or die "can't parse $_";
    my $gn = $1;
    my $draws = $2;
    my @splits = split(/\s*;\s*/, $draws);

    my %rgb = ( 'red' => 0, 'green' => 0, 'blue' => 0);
    
    for my $split (@splits) {
        my @descs = split(/\s*,\s/, $split);
        for my $desc (@descs) {
            $desc =~ /(\d+) (\w+)/;
            my $n = $1;
            my $color = $2;

            $rgb{$color} = max($rgb{$color}, $n);
        }
    }
    my $power = product(values(%rgb));
    $sum += $power;
}

print "$sum\n";
