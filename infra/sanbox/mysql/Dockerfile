FROM mysql:8.0.37

EXPOSE 3306

#MySQL設定ファイルをイメージ内にコピー
ADD my.cnf /etc/mysql/conf.d/my.cnf

#docker runに実行される
CMD ["mysqld"]