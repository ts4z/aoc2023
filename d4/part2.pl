#!/usr/bin/perl -w

# This is not a good answer because it uses brute force.
# We could easily do this in a single pass.

my $sum = 0;
my %cards = ();

sub score {
    my $card = shift;
    my $wins = $card->{wins};
    my $picks = $card->{picks};
    return scalar(grep { ${$card->{wins}}{$_} } (@$picks));
}

for my $line (<ARGV>) {
    $line =~ /Card +(\d+):([\d+ ]+)\s+\|\s+([\d ]+)/ or die "can't parse $line";
    my @wins=grep !/^$/, split(/\s+/,$2);
    my %wins;
    map { $wins{$_} = 1 } (@wins);
    my @picks=grep !/^$/, split(/\s+/, $3);

    $cards{$1} = { number => $1, wins => \%wins, picks => \@picks };
}

my @queue = (values(%cards));
my @cards_held = (values(%cards));

while (my $card = shift @queue) {
    # print "card $card $card->{number}\n";
    my $score = score($card);
    # print "card $card->{number} scores $score\n";
    for (my $i = 1; $i <= $score; $i++) {
        # print "card $card->{number} enqueues ", $i + $card->{number}, "\n";
        my $w = $cards{$i + $card->{number}};
        die "uhoh, can't happen" if (!defined($w));
        push @cards_held, $w;
        push @queue, $w;
    }
}

print "cards won: ", scalar(@cards_held), "\n";
