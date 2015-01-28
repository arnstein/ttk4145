with Ada.Text_IO; use Ada.Text_IO;
with Ada.Strings.Unbounded; use Ada.Strings.Unbounded;
with GNAT.Sockets; use GNAT.Sockets;
with Ada.Streams;
use type Ada.Streams.Stream_Element_Count;

procedure udpListener is
  client  : Socket_Type; Address : Sock_Addr_Type;
  Channel : Stream_Access;

  Message : Unbounded_String;

  Offset : Ada.Streams.Stream_Element_Count;
  Data   : Ada.Streams.Stream_Element_Array (1 .. 256);

begin

  GNAT.Sockets.Initialize;  -- initialize a specific package
  Create_Socket(client,  Family_Inet, Socket_Datagram);

  Address.Port := 30000;
  Address.Addr := Any_Inet_Addr;

  Bind_Socket (client, Address);

  Channel := Stream (client, Address);
  Ada.Text_IO.Put_Line ("Created Stream");

   loop
      Ada.Streams.Read (Channel.All, Data, Offset);
      Address := Get_Address (Channel);
      Ada.Text_IO.Put_Line ("from " & Image (Address) & " : ");
      exit when Offset = 0;
      for I in 1 .. Offset loop
         Ada.Text_IO.Put (Character'Val (Data(I)));
      end loop;
      Put_Line("");
   end loop;

end udpListener;
