class Passman < Formula
    homepage "https://github.com/olishmollie/passman"
    url "https://github.com/olishmollie/passman/archive/v0.5.6.tar.gz"
  
    depends_on "tree"
  
    def install
      bin.install "bin/passman"
      bash_completion.install "completions/passman-completion.bash"
    end
  
  end