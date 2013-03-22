(defun println line (print line "\n"))

(defun - a b (+ a (* -1 b)))

(defun ! a (if a 0 1))

(defun != a b (! (== a b)))

(defun getint (int (getline)))