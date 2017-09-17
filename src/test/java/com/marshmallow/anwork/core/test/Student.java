package com.marshmallow.anwork.core.test;

import com.marshmallow.anwork.core.ProtobufSerializer;
import com.marshmallow.anwork.core.Serializable;
import com.marshmallow.anwork.core.Serializer;
import com.marshmallow.anwork.core.test.protobuf.StudentProtobuf;

import java.io.IOException;

/**
 * This is a dummy object to be used in tests.
 *
 * <p>
 * Created Sep 8, 2017
 * </p>
 *
 * @author Andrew
 */
public class Student implements Serializable<StudentProtobuf> {

  public static final Serializer<Student> SERIALIZER
      = new ProtobufSerializer<StudentProtobuf, Student>(() -> new Student(),
                                                         StudentProtobuf.parser());

  private String name;
  private int id;

  // This constructor is used in the serializer above.
  private Student() {
    this("", 0);
  }

  public Student(String name, int id) {
    this.name = name;
    this.id = id;
  }

  public String getName() {
    return name;
  }

  public int getId() {
    return id;
  }

  @Override
  public StudentProtobuf marshall() throws IOException {
    return StudentProtobuf.newBuilder().setName(getName()).setId(getId()).build();
  }

  @Override
  public void unmarshall(StudentProtobuf t) throws IOException {
    name = t.getName();
    id = t.getId();
  }
}
