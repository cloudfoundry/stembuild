require 'digest'
require 'json'

depObj = {}

@dependencies = JSON.parse(File.read(ARGV[0]))

@dependencies.each do |dep|
  digest = Digest::SHA256.file(dep["file_source"]).hexdigest
  version = File.read(dep["version_source"]).chomp
  depObj[File.basename(dep["file_source"])] = {"sha"=> digest, "version" =>  version}
end

puts depObj.to_json
