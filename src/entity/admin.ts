import { Column, Entity, PrimaryGeneratedColumn } from 'typeorm';

@Entity()
export class Admin {
  @PrimaryGeneratedColumn({
    unsigned: true,
  })
  id?: number;

  @Column('varchar', {
    length: 20,
    unique: true,
    comment: '用户名',
  })
  username: string;

  @Column('text', {
    comment: '密码',
  })
  password: string;

  @Column('varchar', {
    length: 10,
    nullable: true,
    comment: '称呼',
  })
  call?: string;

  @Column('char', {
    length: 11,
    nullable: true,
    comment: '手机号',
  })
  phone?: string;

  @Column('varchar', {
    length: 50,
    nullable: true,
    comment: '电子邮件',
  })
  email?: string;

  @Column('text', {
    nullable: true,
    comment: '头像',
  })
  avatar?: string;

  @Column('tinyint', {
    width: 1,
    unsigned: true,
    default: 1,
    comment: '状态',
  })
  status?: number;

  @Column('int', {
    width: 10,
    unsigned: true,
    default: 0,
    comment: '创建时间',
  })
  create_time: number;

  @Column('int', {
    width: 10,
    unsigned: true,
    default: 0,
    comment: '更新时间',
  })
  update_time: number;
}
