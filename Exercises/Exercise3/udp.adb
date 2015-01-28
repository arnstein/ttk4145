with Ada.Text_IO; use Ada.Text_IO;
with Ada.Strings.Unbounded; use Ada.Strings.Unbounded;
with GNAT.Sockets; use GNAT.Sockets;
with Ada.Streams;
use type Ada.Streams.Stream_Element_Count;

procedure udp is
  client  : Socket_Type; address : Sock_Addr_Type;
  Channel : Stream_Access;

  Message : Unbounded_String;

begin

  GNAT.Sockets.Initialize;  -- initialize a specific package
  Create_Socket(client,  Family_Inet, Socket_Datagram);

  address.Port := 30000;
  address.Addr := Any_Inet_Addr;

  Bind_Socket (client, address);

  Channel := Stream (client, address);
  Ada.Text_IO.Put_Line ("Created Stream");

  loop
    Message := Unbounded_String'Input (Channel);
    Ada.Text_IO.Put_Line ("Recieved input");

    address := Get_Address (Channel);
    Ada.Text_IO.Put_Line (To_String(Message) & " from " & Image (address));

  end loop;
end udp;
