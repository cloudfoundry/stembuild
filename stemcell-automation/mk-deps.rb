 require 'digest'
 require 'json'

*DEP_FILES = ARGV

depObj = {}

DEP_FILES.each do |dep|
  digest = Digest::SHA256.file(dep).hexdigest
  depObj[dep] = {"sha"=> digest, "version" =>  ""}
end

puts depObj.to_json
