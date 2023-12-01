#!/usr/bin/perl -w

use strict;

my $ttl = 0;

sub left {
    my $s = shift;
    if ($s eq '') {
        die 'empty';
    }
    my $ch = substr($s, 0, 1);
    if ($ch =~ /[0-9]/) {
        return $ch;
    }
    if ($s =~ /^one/) { return 1; }
    if ($s =~ /^two/) { return 2; }
    if ($s =~ /^three/) { return 3; }
    if ($s =~ /^four/) { return 4; }
    if ($s =~ /^five/) { return 5; }
    if ($s =~ /^six/) { return 6; }
    if ($s =~ /^seven/) { return 7; }
    if ($s =~ /^eight/) { return 8; }
    if ($s =~ /^nine/) { return 9; }

    return left(substr($s, 1));
}

sub right {
    my $s = shift;
    if ($s eq '') {
        die 'empty';
    }
    my $ch = substr($s, -1, 1);
    print "right $_ ch $ch\n";
    if ($ch =~ /[0-9]/) {
        return $ch;
    }
    if ($s =~ /one$/) { return 1; }
    if ($s =~ /two$/) { return 2; }
    if ($s =~ /three$/) { return 3; }
    if ($s =~ /four$/) { return 4; }
    if ($s =~ /five$/) { return 5; }
    if ($s =~ /six$/) { return 6; }
    if ($s =~ /seven$/) { return 7; }
    if ($s =~ /eight$/) { return 8; }
    if ($s =~ /nine$/) { return 9; }

    return right(substr($s, 0, length($s)-1));
}

while (<ARGV>) {
    chomp;
    my $n = left($_) . right($_);
    print "read $_ -> $n\n";
    $ttl += $n;
}

print "$ttl\n";
