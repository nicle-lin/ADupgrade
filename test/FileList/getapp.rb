apps = Array.new
Dir.foreach(ARGV[0]) { |entry| apps << entry if entry =~ /app\d{1,2}/}
apps.sort! { |x, y| x[3..-1].to_i <=> y[3..-1].to_i }
puts apps
